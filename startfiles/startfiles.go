package startfiles

import (
	"log"
	"os"
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
		err := os.WriteFile("secret.yaml", []byte(`token: "DONATEX TOKEN HERE"`), 0644)
		if err != nil {
			return err
		}
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