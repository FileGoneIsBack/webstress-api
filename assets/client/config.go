package main

import (
	"encoding/json"
	"os"
	"api/core/models/log"
)

var (
	Config = new(conf)
)

type conf struct {
	Name   	 string `json:"name"`
	Slots  	 string `json:"slots"`
	Type   	 string `json:"type"`
	Key	     string `json:"key"`
	Master 	 string `json:"master"`
	MThread  int `json:"maxthreads"`
}

func Load() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(file).Decode(&Config)
	log.Println("succesfully read config! (s.name=\"" + Config.Name + "\")")
}
