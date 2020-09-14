package configs

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	configFileName = "config.yaml"
)

type Config struct {
	Token   string
	IsDebug bool
	Update  tgbotapi.UpdateConfig
}

func LoadConfig() (config Config) {
	config = Config{
		Token:   "",
		IsDebug: true,
	}

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
