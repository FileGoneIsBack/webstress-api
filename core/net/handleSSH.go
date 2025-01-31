package net

import (
    "api/core/models"
	"api/core/database"
	"api/core/net/term"
	"api/core/net/sessions"
	"api/core/net/commands"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Admin struct {
	conn ssh.Session
}

func NewAdmin(ssh ssh.Session) *Admin {
	return &Admin{ssh}
}


func (net *Admin) Handle() {
	//more auth
	net.conn.Write([]byte("\033[?1049h"))
	net.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	defer func() {
		net.conn.Write([]byte("\033[?1049l"))
	}()

    userInfo, err := database.Container.GetUser(net.conn.User())
    if err != nil {
        net.conn.Write([]byte(fmt.Sprintf("Error retrieving user info: %v", err)))
        return
    }
	if !userInfo.HasPermission("cnc") && !userInfo.HasPermission("admin") {
		net.conn.Write([]byte(fmt.Sprint("Must buy CnC!\r\n")))
		time.Sleep(5 * time.Second)
		net.conn.Close()
	}
	netconn := term.NewSSHConn(net.conn)
	net.conn.Write([]byte("\r\n\033[0m"))
	var Session = &sessions.Session{
		User: userInfo,
		Conn: netconn,
	}
	Session.CreateCookie()
	sessions.SessionMutex.Lock()
    for _, session := range sessions.Sessions {
        if session.User.Username == userInfo.Username {
            net.conn.Write([]byte(fmt.Sprintf("Session Already Open!\r\n")))
            logger.Println(session.User.Username + " already has a session open!")
            sessions.SessionMutex.Unlock()
            return
        }
    }
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
			term.SetTitle(net.conn, " ["+SpinnerChars[i%len(SpinnerChars)]+"] "+models.Config.Name+" CnC - Username ["+Session.User.Username+"] - Online ["+strconv.Itoa(sessions.Count())+"] - Rank ["+role+"]")
			i++
		}
		
	}()
	net.conn.Write([]byte("\033[2J\033[1;1H"))

	//handler
	commands.Commands["splash-home"].Exec(Session, nil)
	for {
        term := terminal.NewTerminal(net.conn, userInfo.Username+"@"+models.Config.Name+"$ ")
        cmd, err := term.ReadLine()
        cmd = strings.ToLower(cmd) 
        if err != nil {
            return
        }

        cmdlist := strings.Split(cmd, " ")
        if !commands.IsCommand(cmdlist[0]) {
            fmt.Fprintf(net.conn, "Command (%s) is Invalid\r\n", cmdlist[0])
        } else {
            commands.Commands[cmdlist[0]].Exec(Session, cmdlist)
        }
		if err != nil {
			net.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
		} 
	}
}

