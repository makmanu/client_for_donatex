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

	fmt.Printf("Loaded config: URL=%s, Token=%s, Port=%d\n", cfg.URL, cfg.Token, cfg.Port)

	// Initialize the client
	c := client.NewClient(cfg.URL, cfg.Token)

	// Start the listener
	go listener.StartListener(cfg)

	// Example: Get donations
	err = c.GetDonations(0, 4, "true") // skip 0, take 4
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = c.TestDonations(228, "mrHrunDell", "Проверка донатов)", "RUB", false)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Keep the program running
	select {}
}
