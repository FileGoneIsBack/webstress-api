package servers

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"slices"
	"strconv"
	"strings"
  "encoding/json"
)

func Listen() {
	if Config == nil {
		log.Fatal("Config is not initialized")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", Config.Listener))
	if err != nil {
		log.Fatal(err)
	}
	logger.Println("listening for incoming connections on \"" + fmt.Sprint(Config.Listener) + "\"")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println(err)
			continue
		}
		server := New(conn)
		go server.Handle()
	}
}

func (s *Server) Handle() {
	// Check if the remote address is allowed
	if !slices.Contains(Config.Allowed, s.RemoteAddr()) {
		logger.Println(s.RemoteAddr() + " unauthorized connection!")
		s.NewMessage(MessageFailure, "unauthorized connection!")
		return
	}

	// Read the authentication key
	buf := make([]byte, len(Config.Key))
	n, err := s.Read(buf)
	if err != nil {
		logger.Println("Error while reading authentication key from", s.RemoteAddr(), ":", err)
		return
	}
	if n != len(Config.Key) || buf == nil {
		logger.Println(s.RemoteAddr() + " invalid key length or empty buffer!")
		return
	}
	if !bytes.Equal(buf, []byte(Config.Key)) {
		logger.Println(s.RemoteAddr() + " failed to authenticate!")
		s.NewMessage(MessageFailure, "failed to authenticate!")
		return
	}

	// Read the message after successful authentication
	message, err := s.ReadMessage()
	if err != nil {
		logger.Println(s.RemoteAddr() + " error while reading message: " + err.Error())
		return
	}

	// Check message ID and log mismatches
	if message.ID != MessageAuthenticate {
		logger.Printf("Message ID mismatch from %s: expected=%d, got=%d\n", s.RemoteAddr(), MessageAuthenticate, message.ID)
		return
	}

	// Process the message content
	data := strings.Split(string(message.Content), "|")
	if len(data) < 3 {
		logger.Println("Invalid message format from " + s.RemoteAddr())
		return
	}

	name := data[0]
	slots, err := strconv.Atoi(data[1])
	if err != nil {
		logger.Println("Invalid number of slots from " + s.RemoteAddr() + ": " + data[1])
		return
	}
	stype, err := strconv.Atoi(data[2])
	if err != nil {
		logger.Println("Invalid server type from " + s.RemoteAddr() + ": " + data[2])
		return
	}

	// Update server information
	s.Slots = slots
	s.Type = stype
	s.Name = name
	Servers[s.Name] = s

	// Log registration
	logger.Printf("%s registered as \"%s\" with \"%d\" slots\n", s.RemoteAddr(), name, slots)
	// Send success response
	s.NewMessage(MessageSuccess, "authenticated!")
 
  methods, err := s.ReadMessage()
  if err != nil {
    logger.Println(s.RemoteAddr() + " error while reading message: " + err.Error())
		return
  }
  log.Println(methods)
  
  err = json.Unmarshal(methods.Content, &s.Methods)
  if err != nil {
    logger.Println(s.RemoteAddr() + " failed to parse methods: " + err.Error())
    return
  }
  
  s.NewMessage(MessageSuccess, "synced methods!")
   	logger.Printf("%s synced methods: %v\n", s.RemoteAddr(), s.Methods)

	// Start ongoing operations and keep the connection alive
	go s.Ongoing()
	s.KeepAlive()
}
