package main

import (
	"log"
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
	database.Container.DeleteExpiredInvites()

	bot, err := bot.Init()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	go server.Start()
	handler.HandleUpdates(bot)
}