package startfiles

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func CheckMandatoryFiles() error {
	// Check if config.yaml exists
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		err := os.WriteFile("config.yaml", []byte(`url: "https://donatex.gg/api/v1/"
VTubeStudio:
  port: 14141
  url: "ws://localhost"
minimumDuration: 2
maximumHotkeysPerDonation: 20
defaultHotkeyCoefficient: 2.5`), 0644)
		if err != nil {
			return err
		}
		log.Println("Created default config.yaml")
	}

	// Check if secret.yaml exists
	if _, err := os.Stat("secret.yaml"); os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("DonateX token can be find here: https://donatex.gg/streamer/settings\nEnter your DONATEX TOKEN: ")
		token, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
		token = strings.TrimSpace(token)
		err = os.WriteFile("secret.yaml", []byte(`token: "`+token+`"`), 0644)
		log.Println("Created default secret.yaml")
	}

	// Check if coefficient.yaml exists
	if _, err := os.Stat("coefficient.yaml"); os.IsNotExist(err) {
		err := os.WriteFile("coefficient.yaml", []byte(`coefficients:
  1: 6
  15: 9999`), 0644)
		if err != nil {
			return err
		}
		log.Println("Created default coefficient.yaml")
	}
	return nil
}