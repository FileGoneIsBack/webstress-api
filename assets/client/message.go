package main

import (
	"api/core/models/log"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	ID        int    `json:"mmid"`
	MessageID int    `json:"messageid"`
	Length    int    `json:"length"`
	Content   []byte `json:"mcontent"`
}

type MethodsMessage struct {
	ID        int      `json:"mmid"`
	MessageID int      `json:"messageid"`
	Methods   []string `json:"methods"`
}

type AttackMessage struct {
	ID        int `json:"mmid"`
	MessageID int `json:"messageid"`
	Data      struct {
		User     int    `json:"user"`
		Target   string `json:"target"`
		Port     string `json:"port"`
		Method   string `json:"method"`
		Duration string `json:"duration"`
	} `json:"data"`
	Options struct {
		Threads string `json:"threads"`
		PPS     string `json:"pps"`
		Subnet  string `json:"subnet"`
	} `json:"options"`
}

func NewMessage(ID int, content string) (*Message, []byte) {
	CurrentID++
	m := &Message{
		ID:        ID,
		MessageID: CurrentID,
		Length:    len(content),
		Content:   []byte(content),
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println("error while encoding message!")
	}
	msg := base64.RawStdEncoding.EncodeToString(bytes)
	return m, []byte(msg)
}

func ReadMessage(conn net.Conn) (*Message, error) {
	length := make([]byte, 1)
	n, err := conn.Read(length)
	if err != nil || n != 1 {
		return nil, err
	}
	buf := make([]byte, length[0])
	n, err = conn.Read(buf)
	if err != nil || n < 0 {
		return nil, err
	}
	m := new(Message)
	decoded, err := base64.RawStdEncoding.DecodeString(string(buf[:n]))
	if err != nil {
		log.Println("failed to decode buffer")
		return nil, err
	}
	err = json.Unmarshal(decoded, &m)
	if string(m.Content) == "ping!" {

	} else {
		log.Println("Decoded mcontent:", string(m.Content))
	}
	CurrentID++
	return m, nil
}

func ReadAttack(conn net.Conn) (*AttackMessage, error) {
	length := make([]byte, 1)
	n, err := conn.Read(length)
	if err != nil || n != 1 {
		return nil, err
	}
	buf := make([]byte, length[0])
	n, err = conn.Read(buf)
	if err != nil || n != int(length[0]) {
		return nil, err
	}
	m := new(AttackMessage)
	decoded, err := base64.RawStdEncoding.DecodeString(string(buf[:n]))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(string(decoded))
	//fmt.Println(string(decoded))
	err = json.Unmarshal(decoded, &m)
	if err != nil {

		log.Println(err)
		return nil, err
	}
	CurrentID++
	return m, nil
}

func NewMethodsMessage(ID int, methods []string) (*MethodsMessage, []byte) {
	CurrentID++
	log.Printf("Creating MethodsMessage with methods: %v", methods)

	m := &MethodsMessage{
		ID:        ID,
		MessageID: CurrentID,
		Methods:   methods,
	}

	// Marshal the message to JSON
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println("error while encoding methods message:", err)
		return nil, nil
	}

	// Log the JSON before encoding to base64
	log.Printf("JSON encoded methods message: %s", string(bytes))

	// Encode to base64
	msg := base64.RawStdEncoding.EncodeToString(bytes)

	// Log the base64 message
	log.Printf("Base64 encoded methods message: %s", msg)
	return m, []byte(msg)
}
