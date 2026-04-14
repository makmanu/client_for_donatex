package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/planner"
	"github.com/makmanu/client_for_donatex/plugin"
)

// Start begins the interactive console prompt
func Start(exitChan chan struct{}, c *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	help_text := `=== Client_for_Donatex Console ===
Commands:
  getdonations <skip> <take> <hideTest> - Get donations with pagination and test donation filter
  update                                - Update current hotkeys from VTube Studio
  execute <id|name>                     - Execute a hotkey by ID or name
  executetime <id|name> <seconds>       - Schedule a hotkey to execute for a certain time
  remove <id|name>                      - Remove a hotkey from the schedule
  help                                  - Show this help message
  exit                                  - Exit app
=====================`
	fmt.Println(help_text)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "executetime":
			if len(parts) < 3 {
				fmt.Println("Usage: executetime <id|name> <seconds>")
				continue
			}
			identifierCmd := strings.Join(parts[1:len(parts)-1], " ")
			secondsCmd := parts[len(parts)-1]
			executetimeCmd(identifierCmd, secondsCmd)

		case "getdonations":
			if len(parts) != 4 {
				fmt.Println("Usage: getdonations <skip> <take> <hideTest>")
				continue
			}
			skip, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid skip value, should be an number")
				continue
			}
			take, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid take value, should be an number")
				continue
			}
			if parts[3] != "true" && parts[3] != "false" {
				fmt.Println("Invalid hideTest value, should be 'true' or 'false'")
				continue
			}
			hideTest := parts[3]
			getDonationsCmd(c, skip, take, hideTest)

		case "execute":
			if len(parts) < 2 {
				fmt.Println("Usage: execute <id|name>")
				continue
			}
			identifier := strings.Join(parts[1:], " ")
			executeHotkeyCmd(identifier)

		case "update":
			getCurrentHotkeysCmd()
		
		case "help":
			fmt.Println(help_text)
		
		case "remove":
			if len(parts) < 2 {
				fmt.Println("Usage: remove <id|name>")
				continue
			}
			identifier := strings.Join(parts[1:], " ")
			removeItemFromScheduleCmd(identifier)

		case "exit":
			fmt.Println("Exiting app...")
			exitChan <- struct{}{}
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func executeHotkeyCmd(identifier string) {
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(identifier)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := plugin.ExecuteHotkey(hotkey); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Executed hotkey: %s\n", hotkey.Name)
	}
}

func getCurrentHotkeysCmd() {
	if err := plugin.GetCurrentHotkeys(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved current hotkeys\n")
	}
}

func getDonationsCmd(c *client.Client, skip, take int, hideTest string) {
	if err := c.GetDonations(skip, take, hideTest); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved donations (skip: %d, take: %d)\n", skip, take)
	}
}

func executetimeCmd(key string, time string) {
	time_in_seconds, err := strconv.Atoi(time)
	if err != nil {
		fmt.Println("Invalid time value, should be an number")
		return
	}
	err = planner.DefaultPlanner.AddHotkeyToSchedule(key, time_in_seconds)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Added hotkey %s to schedule to execute for %d seconds\n", key, time_in_seconds)
}

func removeItemFromScheduleCmd(itemName string) {
	err := planner.DefaultPlanner.RemoveItemFromSchedule(itemName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Removed hotkey %s from schedule\n", itemName)
}