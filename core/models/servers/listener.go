package servers

import (
	"api/core/models/log"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
)

func Listen() {
	if Config == nil {
		log.Fatal("Config is not initialized")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", Config.Listener))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening for incoming connections on \"" + fmt.Sprint(Config.Listener) + "\"")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		server := New(conn)
		go server.Handle()
	}
}

func (s *Server) Handle() {
	// Check if the remote address is allowed
	if !slices.Contains(Config.Allowed, s.RemoteAddr()) {
		log.Println(s.RemoteAddr() + " unauthorized connection!")
		s.NewMessage(MessageFailure, "unauthorized connection!")
		return
	}

	// Read the authentication key
	buf := make([]byte, len(Config.Key))
	n, err := s.Read(buf)
	if err != nil {
		log.Println("Error while reading authentication key from", s.RemoteAddr(), ":", err)
		return
	}
	if n != len(Config.Key) || buf == nil {
		log.Println(s.RemoteAddr() + " invalid key length or empty buffer!")
		return
	}
	if !bytes.Equal(buf, []byte(Config.Key)) {
		log.Println(s.RemoteAddr() + " failed to authenticate!")
		s.NewMessage(MessageFailure, "failed to authenticate!")
		return
	}

	// Read the message after successful authentication
	message, err := s.ReadMessage()
	if err != nil {
		log.Println(s.RemoteAddr() + " error while reading message: " + err.Error())
		return
	}

	// Check message ID and log mismatches
	if message.ID != MessageAuthenticate {
		log.Printf("Message ID mismatch from %s: expected=%d, got=%d\n", s.RemoteAddr(), MessageAuthenticate, message.ID)
		return
	}

	// Process the message content
	data := strings.Split(string(message.Content), "|")
	if len(data) < 3 {
		log.Println("Invalid message format from " + s.RemoteAddr())
		return
	}

	name := data[0]
	slots, err := strconv.Atoi(data[1])
	if err != nil {
		log.Println("Invalid number of slots from " + s.RemoteAddr() + ": " + data[1])
		return
	}
	stype, err := strconv.Atoi(data[2])
	if err != nil {
		log.Println("Invalid server type from " + s.RemoteAddr() + ": " + data[2])
		return
	}

	// Update server information
	s.Slots = slots
	s.Type = stype
	s.Name = name
	Servers[s.Name] = s

	// Log registration
	log.Printf("%s registered as \"%s\" with \"%d\" slots\n", s.RemoteAddr(), name, slots)
	// Send success response
	s.NewMessage(MessageSuccess, "authenticated!")

	methods, err := s.ReadMessage()
	if err != nil {
		log.Println(s.RemoteAddr() + " error while reading message: " + err.Error())
		return
	}
	log.Println(methods)

	err = json.Unmarshal(methods.Content, &s.Methods)
	if err != nil {
		log.Println(s.RemoteAddr() + " failed to parse methods: " + err.Error())
		return
	}

	s.NewMessage(MessageSuccess, "synced methods!")
	log.Printf("%s synced methods: %v\n", s.RemoteAddr(), s.Methods)

	// Start ongoing operations and keep the connection alive
	go s.Ongoing()
	s.KeepAlive()
}
