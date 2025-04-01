package main

import (
	"api/core/models/log"
	"encoding/json"
	"os"
)

func LoadMethods() {
	file, err := os.Open("methods.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&methods); err != nil {
		log.Fatal("Failed to decode methods.json:", err)
	}

	log.Println(len(methods), "registered methods!")
}
