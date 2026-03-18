package main

import (
	"fmt"
	"log"

	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
)

func main() {
	fmt.Println("Starting donatex API client...")

	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config: URL=%s, Token=%s\n", cfg.URL, cfg.Token)

	// Initialize the client
	c := client.NewClient(cfg.URL, cfg.Token)

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

	fmt.Println("Client initialized and request made successfully.")
}