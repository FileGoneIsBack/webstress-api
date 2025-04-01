package functions

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
	"api/core/models"
	"github.com/shirou/gopsutil/cpu"
	"api/core/models/log"
)

// commandListener listens for user input in the terminal
func CommandListener() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Command listener started. Type 'exit' to quit.")

	// Listen for commands in a loop
	for {
		if scanner.Scan() {
			cmd := scanner.Text()

			// Process command
			switch cmd {
			case "clear":
				ClearScreen()
				fmt.Println("Screen Cleared!")
			case "reload":
				models.ReloadConfigs()
			case "shutdown":
				fmt.Println("Exiting command listener...")
				return
			case "info":
				ShowSystemInfo()
			default:
				fmt.Println("Unknown command:", cmd)
			}
		}

		// Check for scanner error
		if scanner.Err() != nil {
			log.Println("Error reading input:", scanner.Err())
			return
		}
	}
}

// Function to clear the screen in terminal
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	time.Sleep(5 * time.Millisecond)
}

// Show system info like CPU usage
func ShowSystemInfo() {
	// Get the CPU usage percentage
	percent, err := cpu.Percent(0, true)
	if err != nil {
		log.Println("Error retrieving CPU usage:", err)
		return
	}

	// Display the CPU usage information
	fmt.Println("CPU Usage Info:")
	for i, p := range percent {
		fmt.Printf("CPU %d: %.2f%%\n", i, p)
	}
	
}