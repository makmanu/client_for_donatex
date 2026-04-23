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
	if _, err := os.Stat("listHelp.yaml"); os.IsNotExist(err) {
		err := os.WriteFile("listHelp.yaml", []byte(`# Available types
#   Hotkey - any hotkey from VTS.
#     In Args accepts tags, id. 
#     Tags optional. hotkeys with same tags wont play simultaniosly.
#     id mandatory. its hotkeyid.
#   Tint - change color for selected meshes
#     In Args accepts any amount of Colors, that should contain: listOfMeshes, R, G, B.
#     listOfMeshes accepts []string, with any amount of mesh names that should be colored.
#     R 0-255 of red color.
#     G 0-255 of green color.
#     B 0-255 of blue color.


# Example of list.yaml
Commands:
  Angry:
    Type: "Hotkey"
    Price: "8"
    Args:
      Tags: []
      Id: "1623a123fe135bcc351"
  Cat eyes:
    Type: "Hotkey"
    Price: "3.3"
    Args:
      Tags: ["Right eye", "Left eye"]
      Id: "cbab151356f63161fee35672"
  Red eyes:
    Type: "Tint"
    Price:  "4"
    Args:
      Color1:
        ListOfMeshes: ["ArtMesh1", "ArtMesh2", "ArtMesh63"]
        R: 255
        G: 0
        B: 0
  Bee hair:
    Type: "Tint"
    Price: "7"
    Args:
      Black:
        ListOfMeshes: ["ArtMesh53", "ArtMesh23", "ArtMesh523"]
        R: 0
        G: 0
        B: 0
      Yellow:
        ListOfMeshes: ["ArtMesh246", "ArtMesh37", "ArtMesh262"]
        R: 255
        G: 255
        B: 0`), 0644)
		if err != nil {
			return err
		}
		log.Println("Created default listHelp.yaml")
	}
	return nil
}