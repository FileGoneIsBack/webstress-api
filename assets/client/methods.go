package main

import (
	"encoding/json"
	"os"

)

func LoadMethods() {
	file, err := os.Open("methods.json")
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&methods); err != nil {
		logger.Fatal("Failed to decode methods.json:", err)
	}

	logger.Println(len(methods), "registered methods!")
}

