package models

import (
	"log"
	"os"
	"encoding/json"
)

type config struct {
	BotToken string `json:"bot_token"`
	DbPath   string `json:"db_path"`
	Secure	 bool	`json:"secure"`
	Auth	 string `json:"key"`
	Autobuy struct {
		Key 		string `json:"key"`
		Flat 		string `json:"flat"`
	} `json:"autobuy"`
}
var Config config

func InitConfig() {
	var config config

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config: %v", err)
	}
	Config = config
}