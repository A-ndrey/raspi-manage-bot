package bot

import (
	"context"
	"github.com/A-ndrey/raspi-manage-bot/board"
	"github.com/A-ndrey/raspi-manage-bot/camera"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

func Start(ctx context.Context, config configs.Config, monitoring <-chan db.Measurement) (err error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return
	}

	bot.Debug = config.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updatesChan, err := bot.GetUpdatesChan(config.Update)
	if err != nil {
		return
	}

	notifyOwners(bot, "Bot is running")
	defer notifyOwners(bot, "Bot is shutting down")

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updatesChan:
			handleUpdate(bot, update, config)
		case measurement := <-monitoring:
			handleMonitoring(bot, measurement)
		}
	}
}

func handleMonitoring(bot *tgbotapi.BotAPI, measurement db.Measurement) {
	owners, err := db.GetOwners()
	if err != nil {
		log.Println(err)
		return
	}

	for _, chatID := range owners {
		msg := tgbotapi.NewMessage(chatID, "WARNING\n"+measurement.String())
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, config configs.Config) {
	if update.Message != nil {
		handleMessage(bot, update.Message, config)
	}

	if update.CallbackQuery != nil {
		handleCallbackQuery(bot, update.CallbackQuery)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config configs.Config) {
	chatID := message.Chat.ID

	if message.Command() == "auth" {
		handleAuthCommand(bot, chatID, message.CommandArguments(), config)
	}

	switch message.Text {
	case Temperature:
		msg := tgbotapi.NewMessage(chatID, "Choose environment")
		msg.ReplyMarkup = GetTempInlineKeyboard()
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case Camera:
		msg := tgbotapi.NewMessage(chatID, "Choose mode")
		msg.ReplyMarkup = GetCameraInlineKeyboard()
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case Reboot:
		handleRebootButton(bot, chatID)
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	switch query.Data {
	case BoardTempCallback:
		editMessage(bot, chatID, messageID, "Board temperature:", GetEmptyInlineKeyboard())
		handleInternalTempQuery(bot, chatID)
	case AirTempCallback:
		log.Println(AirTempCallback, "UNIMPLEMENTED")
	case PhotoCallback:
		editMessage(bot, chatID, messageID, "Wait a moment, it may take a few seconds", GetEmptyInlineKeyboard())
		handlePhotoQuery(bot, chatID)
	case VideoCallback:
		editMessage(bot, chatID, messageID, "Choose timings", GetVideoInlineKeyboard())
	case Video5sCallback:
		editMessage(bot, chatID, messageID, "Wait a moment, it may take a few seconds", GetEmptyInlineKeyboard())
		handleVideoQuery(bot, chatID, 5*time.Second)
	case Video10sCallback:
		editMessage(bot, chatID, messageID, "Wait a moment, it may take a few seconds", GetEmptyInlineKeyboard())
		handleVideoQuery(bot, chatID, 10*time.Second)
	case Video15sCallback:
		editMessage(bot, chatID, messageID, "Wait a moment, it may take a few seconds", GetEmptyInlineKeyboard())
		handleVideoQuery(bot, chatID, 15*time.Second)
	}
}

func handleVideoQuery(bot *tgbotapi.BotAPI, chatID int64, duration time.Duration) {
	if !isAuthorized(chatID, db.ROLE_OWNER) {
		sendPermissionDeniedMessage(bot, chatID)
		return
	}

	video, err := camera.TakeVideo(duration)
	if err != nil {
		log.Println(err)
		return
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  "raspi_video",
		Bytes: video,
	}
	msg := tgbotapi.NewVideoUpload(chatID, fileBytes)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func handleInternalTempQuery(bot *tgbotapi.BotAPI, chatID int64) {
	if !isAuthorized(chatID, db.ROLE_OWNER) && !isAuthorized(chatID, db.ROLE_GUEST) {
		sendPermissionDeniedMessage(bot, chatID)
		return
	}

	temp := board.GetTemperature()
	var msg tgbotapi.MessageConfig
	if temp == "" {
		msg = tgbotapi.NewMessage(chatID, "No data about internal temperature")
	} else {
		msg = tgbotapi.NewMessage(chatID, temp)
	}

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func handleRebootButton(bot *tgbotapi.BotAPI, chatID int64) {
	if !isAuthorized(chatID, db.ROLE_OWNER) {
		sendPermissionDeniedMessage(bot, chatID)
		return
	}
	err := board.Reboot()
	if err != nil {
		log.Println(err)
	}
}

func handlePhotoQuery(bot *tgbotapi.BotAPI, chatID int64) {
	if !isAuthorized(chatID, db.ROLE_OWNER) {
		sendPermissionDeniedMessage(bot, chatID)
		return
	}

	photo, err := camera.TakePhoto()
	if err != nil {
		log.Println(err)
		return
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  "raspi_photo",
		Bytes: photo,
	}

	photoMessage := tgbotapi.NewPhotoUpload(chatID, fileBytes)
	_, err = bot.Send(photoMessage)
	if err != nil {
		log.Println(err)
	}
}

func handleAuthCommand(bot *tgbotapi.BotAPI, chatID int64, code string, config configs.Config) {
	if config.Auth.Code != code {
		return
	}

	role := db.ROLE_OWNER

	err := db.InsertAuth(chatID, role)
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Logged in")
	msg.ReplyMarkup = GetOwnerKeyboard()
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func isAuthorized(chatID int64, requiredRole string) bool {
	role := db.GetRoleByChatID(chatID)
	if role == requiredRole {
		return true
	}

	return false
}

func notifyOwners(bot *tgbotapi.BotAPI, text string) {
	chatIDs, err := db.GetOwners()
	if err != nil {
		log.Println(err)
	}

	if chatIDs == nil {
		return
	}

	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func editMessage(
	bot *tgbotapi.BotAPI,
	chatID int64,
	messageID int,
	text string,
	keyboard tgbotapi.InlineKeyboardMarkup,
) {
	msgTxt := tgbotapi.NewEditMessageText(chatID, messageID, text)
	_, err := bot.Send(msgTxt)
	if err != nil {
		log.Println(err)
	}
	msgMrk := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, keyboard)
	_, err = bot.Send(msgMrk)
	if err != nil {
		log.Println(err)
	}
}

func sendPermissionDeniedMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Permission denied")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
