package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/makmanu/client_for_donatex/plugin"
	"github.com/makmanu/client_for_donatex/client"
)

// Start begins the interactive console prompt
func Start(conn *websocket.Conn, exitChan chan struct{}, c *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	help_text := `=== Client_for_Donatex Console ===
Commands:
  createwebhook      - Create a new webhook subscription
  getwebhooks        - List all registered webhooks
  update             - Update current hotkeys from VTube Studio
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
		case "createwebhook":
			CreateWebhookCmd(conn, c)
		case "getwebhooks":
			GetWebhooksCmd(conn, c)
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

func GetWebhooksCmd(conn *websocket.Conn, c *client.Client) {
	webhooks, err := c.GetWebhooks()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, webhook := range webhooks {
		fmt.Printf("✓ Webhook: %+v\n", webhook)
	}
}

func CreateWebhookCmd(conn *websocket.Conn, c *client.Client) {
	webhook, err := c.CreateWebhook()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Created webhook: %+v\n", webhook)
}