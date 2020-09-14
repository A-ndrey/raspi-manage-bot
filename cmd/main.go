package main

import (
	"github.com/A-ndrey/raspi-manage-bot/camera"
	"github.com/A-ndrey/raspi-manage-bot/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	config := configs.LoadConfig()

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = config.IsDebug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updatesChan, err := bot.GetUpdatesChan(config.Update)
	if err != nil {
		log.Fatalln(err)
	}

	for update := range updatesChan {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "photo" {
			handlePhotoCommand(bot, update.Message.Chat.ID)
		}
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
