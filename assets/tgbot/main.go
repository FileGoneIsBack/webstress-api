package main

import (
	"api/core/models/log"
	"bot/database"
 	 handler "bot/core"
	"bot/models/bot"
	"bot/models/server"
	"bot/models"
)

func main() {
	models.InitConfig()
	if err := database.New(); err != nil {
		log.Println("failed to initialize database", err)
		return
	}

	bot, err := bot.Init()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	go server.Start()
	handler.HandleUpdates(bot)
}