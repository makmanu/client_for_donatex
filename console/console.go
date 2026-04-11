package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/makmanu/client_for_donatex/plugin"
)

// Start begins the interactive console prompt
func Start(conn *websocket.Conn, exitChan chan struct{}) {
	reader := bufio.NewReader(os.Stdin)

	help_text := `=== Donatex Console ===
Commands:
  update	         - Update current hotkeys from VTube Studio
  execute <id|name>  - Execute a hotkey by ID or name
  help               - Show this help message
  exit               - Exit console
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
		case "execute":
			if len(parts) < 2 {
				fmt.Println("Usage: execute <id|name>")
				continue
			}
			identifier := strings.Join(parts[1:], " ")
			executeHotkeyCmd(conn, identifier)

		case "update":
			GetCurrentHotkeysCmd(conn)
		
		case "help":
			fmt.Println(help_text)

		case "exit":
			fmt.Println("Exiting app...")
			exitChan <- struct{}{}
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func executeHotkeyCmd(conn *websocket.Conn, identifier string) {
	if err := plugin.ExecuteHotkey(conn, identifier); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Executed hotkey: %s\n", identifier)
	}
}

func GetCurrentHotkeysCmd(conn *websocket.Conn) {
	if err := plugin.GetCurrentHotkeys(conn); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved current hotkeys\n")
	}
}
