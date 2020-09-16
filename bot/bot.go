package bot

import (
	"context"
	"github.com/A-ndrey/raspi-manage-bot/board"
	"github.com/A-ndrey/raspi-manage-bot/camera"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	"github.com/A-ndrey/raspi-manage-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Start(ctx context.Context, config configs.Config) (err error) {
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

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updatesChan:
			handleUpdate(bot, update, config)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, config configs.Config) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "auth":
		handleAuthCommand(bot, chatID, update.Message.CommandArguments(), config)
	case "photo":
		if isAuthorized(bot, chatID, db.ROLE_OWNER) {
			handlePhotoCommand(bot, chatID)
		}
	case "reboot":
		if isAuthorized(bot, chatID, db.ROLE_OWNER) {
			handleRebootCommand(bot, chatID)
		}
	}
}

func handleRebootCommand(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Rebooting board...")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

	err = board.Restart()
	if err != nil {
		log.Println(err)
	}
}

func handlePhotoCommand(bot *tgbotapi.BotAPI, chatID int64) {
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
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func isAuthorized(bot *tgbotapi.BotAPI, chatID int64, requiredRole string) bool {
	role := db.GetRole(chatID)
	if role == requiredRole {
		return true
	}

	msg := tgbotapi.NewMessage(chatID, "You aren't authorized")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	return false
}
