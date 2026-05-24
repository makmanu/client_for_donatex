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
maximumRequestsPerDonation: 20
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
		err = os.WriteFile("secret.yaml", []byte(`token: "`+token+`"
minecraft_server_RCON:
  host: ""
  port: 0
  password: ""`), 0644)
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
    Name: "Angry"
    Type: "Hotkey"
    Price: 8
    Args:
      Id: "349e5f8b62bf43bb8e47fabefd56257a"
  Cat_eyes:
    Name: "Cat eyes"
    Type: "Hotkey"
    Price: 3.3
    Args:
      Id: "bed2c8c430ef41318a5ca63fd6383950"
  Red_eyes:
    Name: "Red eyes"
    Type: "Tint"
    Price:  4
    Args:
      Color1:
        ListOfMeshes: ["ArtMesh1", "ArtMesh2", "ArtMesh63"]
        R: 255
        G: 0
        B: 0
  Bee_hair:
    Name: "Bee hair"
    Type: "Tint"
    Price: 7
    Args:
      Black:
        ListOfMeshes: ["ArtMesh115", "ArtMesh117"]
        R: 0
        G: 0
        B: 0
      Yellow:
        ListOfMeshes: ["ArtMesh114", "ArtMesh123"]
        R: 255
        G: 255
        B: 0
  Black_red_hair:
    Name: "Black red hair"
    Type: "Tint"
    Price: 7
    Args:
      Black:
        ListOfMeshes: ["ArtMesh115", "ArtMesh117"]
        R: 0
        G: 0
        B: 0
      Red:
        ListOfMeshes: ["ArtMesh114", "ArtMesh123"]
        R: 255
        G: 0
        B: 0`), 0644)
		if err != nil {
			return err
		}
		log.Println("Created default listHelp.yaml")
	}

  if _, err := os.Stat("list.yaml"); os.IsNotExist(err) {
		err := os.WriteFile("list.yaml", []byte(`Commands:
  Commands:
  Angry:
    Name: "Angry"
    Type: "Hotkey"
    Price: 8
    Args:
      Id: "349e5f8b62bf43bb8e47fabefd56257a"
  Cat_eyes:
    Name: "Cat eyes"
    Type: "Hotkey"
    Price: 3.3
    Args:
      Id: "bed2c8c430ef41318a5ca63fd6383950"
  Red_eyes:
    Name: "Red eyes"
    Type: "Tint"
    Price:  4
    Args:
      Color1:
        ListOfMeshes: ["ArtMesh1", "ArtMesh2", "ArtMesh63"]
        R: 255
        G: 0
        B: 0
  Bee_hair:
    Name: "Bee hair"
    Type: "Tint"
    Price: 7
    Args:
      Black:
        ListOfMeshes: ["ArtMesh115", "ArtMesh117"]
        R: 0
        G: 0
        B: 0
      Yellow:
        ListOfMeshes: ["ArtMesh114", "ArtMesh123"]
        R: 255
        G: 255
        B: 0
  Black_red_hair:
    Name: "Black red hair"
    Type: "Tint"
    Price: 7
    Args:
      Black:
        ListOfMeshes: ["ArtMesh115", "ArtMesh117"]
        R: 0
        G: 0
        B: 0
      Red:
        ListOfMeshes: ["ArtMesh114", "ArtMesh123"]
        R: 255
        G: 0
        B: 0`), 0644)
    if err != nil {
      return err
    }
    log.Println("Created default list.yaml")
  }
	return nil
}