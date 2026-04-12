package console

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/plugin"
)

// Start begins the interactive console prompt
func Start(exitChan chan struct{}, c *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	help_text := `=== Client_for_Donatex Console ===
Commands:
  deletewebhook <id>                    - Delete a webhook subscription by ID
  createwebhook                         - Create a new webhook subscription
  getwebhooks                           - List all registered webhooks
  getdonations <skip> <take> <hideTest> - Get donations with pagination and test donation filter
  update                                - Update current hotkeys from VTube Studio
  execute <id|name>                     - Execute a hotkey by ID or name
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
		case "createwebhook":
			createWebhookCmd(c)

		case "getwebhooks":
			getWebhooksCmd(c)

		case "deletewebhook":
			if len(parts) < 2 {
				fmt.Println("Usage: deletewebhook <id>")
				continue
			}
			deleteWebhookCmd(c, parts[1])

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
		
		case "Devsign":
			if len(parts) < 3 {
				fmt.Println("Usage: Devsign <secret> <body>")
				continue
			}
			secret := parts[1]
			body := strings.Join(parts[2:], " ")
			convertBodyToSignatureCmd(secret, []byte(body))
		
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
	if err := plugin.ExecuteHotkey(identifier); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Executed hotkey: %s\n", identifier)
	}
}

func getCurrentHotkeysCmd() {
	if err := plugin.GetCurrentHotkeys(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved current hotkeys\n")
	}
}

func getWebhooksCmd(c *client.Client) {
	webhooks, err := c.GetWebhooks()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, webhook := range webhooks {
		fmt.Printf("✓ Webhook: %+v\n", webhook)
	}
}

func createWebhookCmd(c *client.Client) {
	webhook, err := c.CreateWebhook()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Created webhook: %+v\n", webhook)
}

func deleteWebhookCmd(c *client.Client, webhookId string) {
	if err := c.DeleteWebhook(webhookId); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Deleted webhook with ID: %s\n", webhookId)
	}
}

func getDonationsCmd(c *client.Client, skip, take int, hideTest string) {
	if err := c.GetDonations(skip, take, hideTest); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved donations (skip: %d, take: %d)\n", skip, take)
	}
}

func convertBodyToSignatureCmd(secret string, body []byte){
	expectedSignature := hmac.New(sha256.New, []byte(secret))
	expectedSignature.Write(body)
	expectedSignatureHex := hex.EncodeToString(expectedSignature.Sum(nil))
	fmt.Printf("Expected signature: %s\n", expectedSignatureHex)
}