package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"  
	DarkGreen = "\033[32;1m" 
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Pink      = "\033[35;1m" 
	White     = "\033[37m"
)

// Logger instance
var Logger = log.New(os.Stdout, "", log.LstdFlags)

// getColor returns the corresponding color for log levels
func getColor(msg string) string {
	parts := strings.SplitN(msg, " ", 2)
	if len(parts) > 1 {
		msg = parts[1] // Ignore the timestamp
	}
	switch {
	case strings.HasPrefix(msg, "error"):
		return Red
	case strings.HasPrefix(msg, "warn"):
		return Yellow
	case strings.HasPrefix(msg, "info"):
		return Green
	case strings.HasPrefix(msg, "debug"):
		return Blue
	case strings.HasPrefix(msg, "succes"):
		return Green
	case strings.HasPrefix(msg, "[NET]"):
		return Purple
	case strings.HasPrefix(msg, "[SSH]"):
		return Pink
	case strings.HasPrefix(msg, "listening"):
		return DarkGreen
	case strings.HasPrefix(msg, "http"):
		return Purple
	default:
		return White
	}
}


func Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	color := getColor(msg)
	fmt.Printf("%s%s%s\n", color, msg, Reset) // Use fmt.Printf instead of Logger.Printf
}

func Println(v ...interface{}) {
	msg := fmt.Sprint(v...)
	color := getColor(msg)
	fmt.Println(color + msg + Reset) // Use fmt.Println instead of Logger.Println
}

func Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	color := Red
	Logger.Println(color + msg + Reset)
	Logger.Println("force quitting..")
	time.Sleep(200 * time.Millisecond)
	os.Exit(1)
}

func Fatalf(v ...interface{}) {
	msg := fmt.Sprint(v...)
	color := Red
	Logger.Printf("%s%s%s", color, msg, Reset)
	Logger.Println("force quitting..")
	time.Sleep(200 * time.Millisecond)
	os.Exit(1)
}
