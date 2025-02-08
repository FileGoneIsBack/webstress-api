package main

import (
	"log"
	"bot/database"
	"os"
	handler "bot/core"
	"bot/models/bot"
	"bot/models/server"
	"bot/models"
)

var logger = log.New(os.Stderr, "[main] ", log.Ltime|log.Lshortfile)

func main() {
	models.InitConfig()
	if err := database.New(); err != nil {
		logger.Println("failed to initialize database", err)
		return
	}

	bot, err := bot.Init()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	go server.Start()
	handler.HandleUpdates(bot)
}