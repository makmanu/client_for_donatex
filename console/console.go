package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/makmanu/client_for_donatex/RCON_minecraft"
	"github.com/makmanu/client_for_donatex/client"
	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/planner"
	"github.com/makmanu/client_for_donatex/plugin"
)

// Start begins the interactive console prompt
func Start(exitChan chan struct{}, c *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	help_text := `=== Client_for_Donatex Console ===
Commands:
  selectmeshes                                    - return selected meshes in VTubeStudio
  tint <r> <g> <b> <a>                            - Tint model with specified RGBA color
  tintmeshes <r> <g> <b> <a> <mesh1> <mesh2>...   - Tint, but for choosed meshes
  tintmeshesfadein <use command to know>          - Tintmeshes, but with animation 
  rconrepeat <times> <command>                    - Send RCON command to Minecraft server multiple times
  rcon <command>                                  - Send RCON command to Minecraft server
  reqmeshes                                       - Get info about curent model meshes
  getdonations <skip> <take> <hideTest>           - Get donations with pagination and test donation filter
  update                                          - Update current hotkeys from VTubeStudio
  execute <id>                                    - Execute a hotkey by ID or name
  executetime <id> <seconds>                      - Schedule a hotkey to execute for a certain time
  remove <name>                                   - Remove a hotkey from the schedule
  help                                            - Show this help message
  exit                                            - Exit app
=====================`
	fmt.Println(help_text)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "executetime":
			if len(parts) != 3 {
				fmt.Println("Usage: executetime <id> <seconds>")
				continue
			}
			identifierCmd := parts[1]
			seconds := parts[2]
			executetimeCmd(identifierCmd, seconds)

		case "reqmeshes":
			reqMeshesCMD()

		case "selectmeshes":
			selectMeshesCMD()

		case "tintmeshes":
			if len(parts) < 4 {
				fmt.Println("Usage: tintmeshes <r> <g> <b> <mesh1> <mesh2>...")
				continue
			}
			r, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid r value, should be an integer")
				continue
			}
			g, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid g value, should be an integer")
				continue
			}
			b, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Println("Invalid b value, should be an integer")
				continue
			}

			meshes := parts[5:]
			tintMeshesCMD(r, g, b, meshes)

		case "tintmeshesfadein":
			if len(parts) < 9 {
				fmt.Println("Usage: tintmeshesfadein <r> <g> <b> <rnew> <gnew> <bnew> <mesh1> <mesh2>...")
				continue
			}
			r, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid r value, should be an integer")
				continue
			}
			g, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid g value, should be an integer")
				continue
			}
			b, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Println("Invalid b value, should be an integer")
				continue
			}
			rnew, err := strconv.Atoi(parts[5])
			if err != nil {
				fmt.Println("Invalid r value, should be an integer")
				continue
			}
			gnew, err := strconv.Atoi(parts[6])
			if err != nil {
				fmt.Println("Invalid g value, should be an integer")
				continue
			}
			bnew, err := strconv.Atoi(parts[7])
			if err != nil {
				fmt.Println("Invalid b value, should be an integer")
				continue
			}

			meshes := parts[9:]
			go tintMeshesFadeInCMD(r, g, b, rnew, gnew, bnew, meshes)

		case "getdonations":
			if len(parts) != 4 {
				fmt.Println("Usage: getdonations <skip> <take> <hideTest>")
				continue
			}
			skip, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid skip value, should be an number")
				continue
			}
			take, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid take value, should be an number")
				continue
			}
			if parts[3] != "true" && parts[3] != "false" {
				fmt.Println("Invalid hideTest value, should be 'true' or 'false'")
				continue
			}
			hideTest := parts[3]
			getDonationsCmd(c, skip, take, hideTest)

		case "execute":
			if len(parts) != 2 {
				fmt.Println("Usage: execute <id>")
				continue
			}
			identifier := parts[1]
			executeCmd(identifier)

		case "rcon":
			if len(parts) < 2 {
				fmt.Println("Usage: rcon <command>")
				continue
			}
			command := strings.Join(parts[1:], " ")
			rconCmd(command)

		case "rconrepeat":
			if len(parts) < 3 {
				fmt.Println("Usage: rconrepeat <repeat> <command>")
				continue
			}
			repeat, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid repeat value, should be an number")
				continue
			}
			command := strings.Join(parts[2:], " ")
			rconrepeatCmd(command, repeat)

		case "tint":
			if len(parts) != 5 {
				fmt.Println("Usage: tintmodel <r> <g> <b> <a>")
				continue
			}
			r, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid r value, should be an integer")
				continue
			}
			g, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid g value, should be an integer")
				continue
			}
			b, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Println("Invalid b value, should be an integer")
				continue
			}
			a, err := strconv.Atoi(parts[4])
			if err != nil {
				fmt.Println("Invalid a value, should be an integer")
				continue
			}
			tintModelCmd(r, g, b, a)

		case "update":
			getCurrentHotkeysCmd()
		
		case "help":
			fmt.Println(help_text)
		
		case "remove":
			if len(parts) < 2 {
				fmt.Println("Usage: remove <id|name>")
				continue
			}
			identifier := strings.Join(parts[1:], " ")
			removeItemFromScheduleCmd(identifier)

		case "exit":
			fmt.Println("Exiting app...")
			exitChan <- struct{}{}
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func executeCmd(identifier string) {
	command, ok := config.CommandList.Commands[identifier]
	if !ok {
		fmt.Printf("Unknown command identifier '%s', can't execute\n", identifier)
		return
	}
	switch command.Type {
	case "Hotkey":
		hotkeyId := command.Args.Id
		if err := plugin.ExecuteHotkey(hotkeyId); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("✓ Executed hotkey: %s\n", hotkeyId)
		}
	case "Tint":
		for colorName, color := range command.Args.Colors {
			err := plugin.TintMeshesFadeIn(255, 255, 255, color.R, color.G, color.B, 1, color.ListOfMeshes)
			if err != nil {
				fmt.Printf("Error executing tint for color '%s': %v\n", colorName, err)
			} else {
				fmt.Printf("✓ Executed tint command for color '%s'\n", colorName)
			}
		}
	default:
		fmt.Printf("Command with id '%s' has unknown type '%s', can't execute\n", identifier, command.Type)
		return
	}
}

func getCurrentHotkeysCmd() {
	if err := plugin.GetCurrentHotkeys(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved current hotkeys\n")
	}
}

func getDonationsCmd(c *client.Client, skip, take int, hideTest string) {
	if err := c.GetDonations(skip, take, hideTest); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved donations (skip: %d, take: %d)\n", skip, take)
	}
}

func executetimeCmd(identifier string, time string) {
	time_in_seconds, err := strconv.Atoi(time)
	if err != nil {
		fmt.Println("Invalid time value, should be an number")
		return
	}

	command, ok := config.CommandList.Commands[identifier]
	if  !ok {
		fmt.Printf("Command with id '%s' not found\n", identifier)
		return
	}
	switch command.Type {
	case "Hotkey":
		err = planner.DefaultPlanner.AddHotkeyToSchedule(identifier, command.Args.Id, time_in_seconds)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	case "Tint":
		for colorName, color := range command.Args.Colors {
			err = planner.DefaultPlanner.AddTintToSchedule(colorName, color.R, color.G, color.B, time_in_seconds, color.ListOfMeshes)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
	default:
		fmt.Printf("Command with id '%s' has unknown type '%s', can't be scheduled\n", identifier, command.Type)
		return
	}

	fmt.Printf("✓ Added hotkey %s to schedule to execute for %d seconds\n", identifier, time_in_seconds)
}

func removeItemFromScheduleCmd(itemName string) {
	err := planner.DefaultPlanner.RemoveItemFromSchedule(itemName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Removed hotkey %s from schedule\n", itemName)
}

func tintModelCmd(r, g, b, a int) {
	err := plugin.TintModel(r, g, b, a)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Tinted model with color RGBA(%d, %d, %d, %d)\n", r, g, b, a)
}

func reqMeshesCMD() {
	err := plugin.RequestArtMeshList()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func tintMeshesCMD(r, g, b int, meshes []string){
	err := plugin.TintMeshes(r, g, b, meshes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Tinted some meshes with color RGBA(%d, %d, %d, %d)\n", r, g, b, 255)
}

func tintMeshesFadeInCMD(R, G, B, Rnew, Gnew, Bnew int, meshes []string){
	err := plugin.TintMeshesFadeIn(R, G, B, Rnew, Gnew, Bnew, 1, meshes)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✓ Tinted meshes with animation and color RGBA(%d, %d, %d, %d)\n", R, G, B, 255)
}

func selectMeshesCMD() {
	meshes, err := plugin.AskUserForMeshes()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for i, mesh := range meshes {
		meshes[i] = fmt.Sprintf("\"%s\"", mesh)
	}

	fmt.Printf("✓ Selected meshes: [%s]\n", strings.Join(meshes, ", "))
}

func rconCmd(command string) {
	response, err := RCON_minecraft.SendRCONCommand(command)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("RCON response: %s\n", response)
}

func rconrepeatCmd(command string, repeat int) {
	err := RCON_minecraft.SendRCONCommandRepeatedly(command, repeat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Sent RCON command '%s' %d times\n", command, repeat)
}