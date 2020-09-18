package configs

import (
	"github.com/A-ndrey/raspi-manage-bot/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	configFileName = "config.yaml"
)

type Config struct {
	Token string
	Auth  struct {
		Code string
	}
	Debug      bool
	Update     tgbotapi.UpdateConfig
	Monitoring []db.Measurement
}

func LoadConfig() (config Config) {
	config = Config{Debug: true}

	fileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Println("Can't read config.yaml file. Config will have default values.")
		log.Println(err)
		return
	}

	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Println("Can't unmarshal config.yaml file. Config will have default values.")
		log.Println(err)
		return
	}

	return
}
