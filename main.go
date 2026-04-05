package main

import (
	"fmt"
	"log"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/listener"
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

	// Start the listener
	go listener.StartListener(cfg)

	// Example: Get donations
	err = c.GetDonations(0, 4, "true") // skip 0, take 4
	if err != nil {
		fmt.Println("Error:", err)
		return
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
