package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/makmanu/client_for_donatex/RCON_minecraft"
	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/console"
	"github.com/makmanu/client_for_donatex/planner"
	"github.com/makmanu/client_for_donatex/plugin"
	"github.com/makmanu/client_for_donatex/startfiles"
)

func main() {
	// Set up logging to file
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	err = startfiles.CheckMandatoryFiles()
	if err != nil {
		log.Fatalf("Failed to create start files: %v", err)
	}

	log.Println("Starting donatex API client...")

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)	
	}

	secret, err := config.LoadSecret("secret.yaml")
	if err != nil {
		log.Fatalf("Failed to load secret: %v", err)
	}

	if len(secret.Token) < 250 {
		log.Println("Invalid donatex token")
		err := os.Remove("secret.yaml")
		if err != nil {
			log.Printf("Failed to remove invalid secret file: %v", err)
		}
		fmt.Println("DonateX token seems invalid, restart app and enter correct token\nDonateX token can be find here: https://donatex.gg/streamer/settings\nApp will be closed in 5 seconds...")
		time.Sleep(5 * time.Second)
		return
	}

	if cfg.MinecraftServerRCONEnabled == true {
		if err := RCON_minecraft.SetRCONConfiguration(); err != nil {
			fmt.Printf("Failed to set RCON configuration: %v\nRCON disabled for current session\n", err)
			cfg.MinecraftServerRCONEnabled = false
		}
		err := RCON_minecraft.EstablishRCONConnection()
		if err != nil {
			fmt.Printf("Failed to establish RCON connection: %v\nRCON disabled for current session\n", err)
			cfg.MinecraftServerRCONEnabled = false
		} else {
			defer RCON_minecraft.RCONConnection.Close()
			fmt.Println("Successfully established RCON connection to Minecraft server")
		}
	} else {
		fmt.Println("Minecraft server RCON integration is disabled in config, skipping RCON setup")
	}

	log.Printf("Config loaded")

	commandList, err := config.LoadCommandList("list.yaml")
	if err != nil {
		log.Fatalf("Failed to load command list: %v", err)
	}
	log.Printf("Command list loaded")
	log.Println(commandList)

	c := client.NewClient(cfg.URL, secret.Token)

	err = plugin.ConnectPluginWebsocket(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to VTube Studio plugin websocket: %v", err)
	}
	defer plugin.ClosePluginWebsocket()
	log.Println("Connected to VTube Studio plugin websocket")

	err = plugin.SessionAuthPlugin()
	if err != nil {
		log.Fatalf("Failed to authenticate with VTube Studio plugin: %v", err)
	}
	log.Println("Authenticated with VTube Studio plugin")

	log.Println("Getting current hotkeys from VTube Studio...")
	err = plugin.GetCurrentHotkeys()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	} else {
		log.Println("Successfully retrieved current hotkeys from VTube Studio and saved to file plugin/hotkeys.yaml")
	}

	log.Print("Starting planner...\n")
	planner := planner.NewPlanner(cfg.MinimumDuration, cfg.MaximumRequestsPerDonation)
	go planner.Start()
	log.Print("Planner started\n")

	log.Print("Starting SignalR client...\n")
	ctx := context.Background()

	signalrClient, err := client.ConnectWithTokenAutoReconnect(
		ctx,
		"https://donatex.gg",
		secret.Token,
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = signalrClient

	log.Println("signalrClient started")

	// Start the console
	exitChan := make(chan struct{})
	go console.Start(exitChan, c)
	// wait for exit signal
	<-exitChan
	log.Println("Received exit signal, shutting down...")
}
