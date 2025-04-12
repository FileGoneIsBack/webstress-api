package models

import (
	"log"
	"os"
	"encoding/json"
	"strings"
	"math/rand"
	"fmt"
	"time"
)

type config struct {
	BotToken string `json:"bot_token"`
	DbPath   string `json:"db_path"`
	Secure	 bool	`json:"secure"`
	Auth	 string `json:"key"`
	Custom	 string `json:"custom"`
	Database struct {
		Host 			string `json:"host"`
		Database 		string `json:"database"`
		Username 	    string `json:"username"`
		Password 		string `json:"password"`
	} `json:"database"`
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

func GenCode() string {
	name := Config.Custom

	// Replace spaces with dashes
	safeName := strings.ReplaceAll(name, " ", "-")

	// Generate a 4-digit random number
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000 // guarantees 4 digits (1000–9999)

	return fmt.Sprintf("%s-%d", safeName, code)
}