package RCON_minecraft

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	"github.com/makmanu/client_for_donatex/config"
	"gopkg.in/yaml.v3"
)

var RCONConnection *rcon.Conn

func EstablishRCONConnection() (error) {
	if config.Config.MinecraftServerRCONEnabled == false {
		return fmt.Errorf("RCON is disabled")
	}
	
	if RCONConnection != nil {
		RCONConnection.Close()
		RCONConnection = nil
	}
	var err error
	RCONConnection, err = rcon.Dial(fmt.Sprintf("%s:%d", config.Secret.MinecraftRCON.Host, config.Secret.MinecraftRCON.Port), config.Secret.MinecraftRCON.Password)
	return err
}

func SetRCONConfiguration() error {
	if config.Config.MinecraftServerRCONEnabled == false {
		return fmt.Errorf("RCON is disabled")
	}
	if config.Secret.MinecraftRCON.Host == "" || config.Secret.MinecraftRCON.Port == 0 || config.Secret.MinecraftRCON.Password == "" {
		fmt.Println("Enter Minecraft server RCON configuration\nHost: ")
		var host string
		fmt.Scanln(&host)
		if host == "" {
			return fmt.Errorf("host cannot be empty")
		}

		fmt.Println("Port: ")
		var port int
		fmt.Scanln(&port)
		if port <= 0 {
			return fmt.Errorf("port must be a positive integer")
		}

		fmt.Println("Password: ")
		var password string
		fmt.Scanln(&password)
		if password == "" {
			return fmt.Errorf("password cannot be empty")
		}
		
		config.Secret.MinecraftRCON.Host = host
		config.Secret.MinecraftRCON.Port = port
		config.Secret.MinecraftRCON.Password = password
		newSecret, err := yaml.Marshal(config.Secret)
		if err != nil {
			return err
		}
		err = os.WriteFile("secret.yaml", newSecret, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendRCONCommand(command string) (string, error) {
	if config.Config.MinecraftServerRCONEnabled == false {
		return "", fmt.Errorf("RCON is disabled")
	}

	if RCONConnection == nil {
		fmt.Println("RCON connection is not established, trying to establish...")
		if err := EstablishRCONConnection(); err != nil {
			return "", fmt.Errorf("RCON connection is not established: %v", err)
		}
	}

	response, err := RCONConnection.Execute(command)
	if err != nil {
		fmt.Println("RCON execute failed, reconnecting:", err)
		if err2 := EstablishRCONConnection(); err2 != nil {
			return "", fmt.Errorf("RCON execute failed: %v; reconnect failed: %v", err, err2)
		} else {
			fmt.Println("Reconnected to RCON, retrying command...")
		}
		response, err = RCONConnection.Execute(command)
		if err != nil {
			return "", err
		}
	}
	return response, nil
}

func PrepareRCONCommand(command string, args []string) string {
	for i, arg := range args {
		command = strings.Replace(command, fmt.Sprintf("<%d>", i), arg, 1)
	}
	return command
}

func SendRCONCommandRepeatedly(command string, times int) error {
	for range times {
		respond, err := SendRCONCommand(command)
		if err != nil {
			return err
		}
		log.Printf("RCON command response: %s", respond)
		time.Sleep(550 * time.Millisecond)
	}
	return nil
}