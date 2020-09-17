package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	Camera      = "Camera"
	Reboot      = "Reboot"
	Photo       = "Photo"
	Video       = "Video"
	Video5s     = "5s"
	Video10s    = "10s"
	Video15s    = "15s"
	Temperature = "Temperature"
	Board       = "Board"
	Air         = "Air"
	Empty       = ""

	BoardTempCallback = "BoardTempCallback"
	AirTempCallback   = "AirTempCallback"
	PhotoCallback     = "PhotoCallback"
	VideoCallback     = "VideoCallback"
	Video5sCallback   = "Video5sCallback"
	Video10sCallback  = "Video10sCallback"
	Video15sCallback  = "Video15sCallback"
)

var cameraBtn = tgbotapi.NewKeyboardButton(Camera)
var tempBtn = tgbotapi.NewKeyboardButton(Temperature)
var rebootBtn = tgbotapi.NewKeyboardButton(Reboot)

var firstRowOwnerKB = tgbotapi.NewKeyboardButtonRow(cameraBtn, tempBtn, rebootBtn)

var ownerKeyboard = tgbotapi.NewReplyKeyboard(firstRowOwnerKB)

var boardTempBtn = tgbotapi.NewInlineKeyboardButtonData(Board, BoardTempCallback)
var airTempBtn = tgbotapi.NewInlineKeyboardButtonData(Air, AirTempCallback)
var photoBtn = tgbotapi.NewInlineKeyboardButtonData(Photo, PhotoCallback)
var videoBtn = tgbotapi.NewInlineKeyboardButtonData(Video, VideoCallback)
var video5sBtn = tgbotapi.NewInlineKeyboardButtonData(Video5s, Video5sCallback)
var video10sBtn = tgbotapi.NewInlineKeyboardButtonData(Video10s, Video10sCallback)
var video15sBtn = tgbotapi.NewInlineKeyboardButtonData(Video15s, Video15sCallback)
var emptyBtn = tgbotapi.NewInlineKeyboardButtonData(Empty, "")

var rowTempKB = tgbotapi.NewInlineKeyboardRow(boardTempBtn, airTempBtn)
var rowCameraKB = tgbotapi.NewInlineKeyboardRow(photoBtn, videoBtn)
var rowVideoKB = tgbotapi.NewInlineKeyboardRow(video5sBtn, video10sBtn, video15sBtn)
var rowEmptyKB = tgbotapi.NewInlineKeyboardRow(emptyBtn)

var tempKeyboard = tgbotapi.NewInlineKeyboardMarkup(rowTempKB)
var cameraKeyboard = tgbotapi.NewInlineKeyboardMarkup(rowCameraKB)
var videoKeyboard = tgbotapi.NewInlineKeyboardMarkup(rowVideoKB)
var emptyKeyboard = tgbotapi.NewInlineKeyboardMarkup(rowEmptyKB)

func GetOwnerKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return ownerKeyboard
}

func GetTempInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tempKeyboard
}

func GetCameraInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return cameraKeyboard
}

func GetVideoInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return videoKeyboard
}

func GetEmptyInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return emptyKeyboard
}
