package main

import (
	"fmt"
	"log"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/listener"
	"github.com/makmanu/client_for_donatex/plugin"
)

func main() {
	fmt.Println("Starting donatex API client...")

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	secret, err := config.LoadSecret("secret.yaml")
	if err != nil {
		log.Fatalf("Failed to load secret: %v", err)
	}

	fmt.Printf("Loaded config: URL=%s, Token=%s, Port=%d\n", cfg.URL, secret.Token, cfg.Port)

	// Initialize the client
	c := client.NewClient(cfg.URL, secret.Token)

	pluginConn, err := plugin.ConnectPluginWebsocket(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to VTube Studio plugin websocket: %v", err)
	}
	defer pluginConn.Close()
	fmt.Println("Connected to VTube Studio plugin websocket")

	err = plugin.SessionAuthPlugin(pluginConn)
	if err != nil {
		log.Fatalf("Failed to authenticate with VTube Studio plugin: %v", err)
	}
	fmt.Println("Authenticated with VTube Studio plugin")
	
	// Start the listener
	go listener.StartListener(cfg)

	// Example: Get donations
	err = c.GetDonations(0, 4, "true") // skip 0, take 4
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Getting current hotkeys from VTube Studio...")
	err = plugin.GetCurrentHotkeys(pluginConn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Successfully retrieved current hotkeys from VTube Studio and saved to file plugin/hotkeys.yaml")
	}

	/*err = c.TestDonations(228, "mrHrunDell", "Проверка донатов)", "RUB", false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}*/

	/*webhook, err := c.CreateWebhook("https://makmanu.com:3000/webhook", "DonationCreated", "Client_1", secret.WebhookSecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created webhook: %+v\n", webhook)*/

	// Keep the program running
	select {}
}
