package commands

import (
    "fmt"
    "api/core/net/sessions"
	"api/core/database/atks"
)

func ongoing(session *sessions.Session, args []string) {
    attacks, err := atks.Container.GetRunning(session.User)
    if err != nil {
        fmt.Fprintf(session.Conn, "Error: %v\n\r", err)
        return
    }
	
	fmt.Fprintf(session.Conn, "\033c")
    if len(attacks) > 0 {
        fmt.Fprintf(session.Conn, "+------------------------------------------------------------------------------+\n\r")
        fmt.Fprintf(session.Conn, "|     ID     |         Host         |       Duration       |      Method       |\n\r")
        fmt.Fprintf(session.Conn, "+------------------------------------------------------------------------------+\n\r")

        for _, attack := range attacks {
            fmt.Fprintf(session.Conn, "| %-10d | %-20s | %-20d | %-17s |\n\r",
                attack.ID, attack.Target, attack.Duration, attack.Method.Name)
        }

        fmt.Fprintf(session.Conn, "+------------------------------------------------------------------------------+\n\r")
    } else {
        fmt.Fprintf(session.Conn, "No running attacks\n\r")
    }
}