package bot

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"bot/models"
)

func Init() (*tgbotapi.BotAPI, error) {
	log.Println(models.Config.BotToken)
	bot, err := tgbotapi.NewBotAPI(models.Config.BotToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}