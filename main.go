package main

import (
	"fmt"
	"log"
	"os"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/console"
	"github.com/makmanu/client_for_donatex/listener"
	"github.com/makmanu/client_for_donatex/plugin"
)

func main() {
	// Set up logging to file
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

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
	go listener.StartListener(cfg, secret)

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

	/*fmt.Println("Trying to execute 1 hotkey")
	err = plugin.ExecuteHotkey(pluginConn, "1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Successfully activated hotkey with ID 1")
	}*/

	/*fmt.Println("Trying to execute 1 hotkey in 2 seconds")
	go func() {
		time.Sleep(2 * time.Second)
		err = plugin.ExecuteHotkey(pluginConn, "1")
		if err != nil {
			fmt.Println("Error:", err)
			return
		} else {
			fmt.Println("Successfully activated hotkey with ID 1")
		}
	}()*/

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

	// Start the console
	exitChan := make(chan struct{})
	go console.Start(pluginConn, exitChan)
	// wait for exit signal
	select {
	case <-exitChan:
		fmt.Println("Received exit signal, shutting down...")
		return
	}
}
