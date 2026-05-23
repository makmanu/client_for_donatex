package planner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/makmanu/client_for_donatex/config"
	"github.com/makmanu/client_for_donatex/structs"
)

func GetRequestsFromMessage(message string) (map[string][]string, error) {
	result := make(map[string][]string)
	
	re := regexp.MustCompile(`&([^#]+)#([0-9]+)`)
	matches := re.FindAllStringSubmatch(message, -1)

	for _, m := range matches {
		identifier := m[1]
		args := strings.Split(m[2], ";")

		result[identifier] = args
	}

	return result, nil
}

func HandleDonation(donation structs.Donation) {
	requests, err := GetRequestsFromMessage(donation.Message)
		if err != nil {
			fmt.Printf("Error occurred while parsing hotkeys: %v\n", err)
			return
		}
		if len(requests) > DefaultPlanner.maximumRequestsPerDonation {
			fmt.Printf("Too many requests: %d. Maximum is %d(can be changed in config), skipping...\n", len(requests), DefaultPlanner.maximumRequestsPerDonation)
			return
		}
		if len(requests) > 0 {
			donationAmount := 0.0
			if donation.Currency == "RUB" {
				donationAmount = donation.Amount
			} else {
				donationAmount = donation.AmountInRub
			}
			fmt.Printf("? Donation amount: %.2f RUB\n", donationAmount)
			for identifier, args := range requests {
				donationAmount = HandleRequest(identifier, args, donationAmount)
				fmt.Printf("Remaining amount after handling request '%s': %.2f RUB\n", identifier, donationAmount)
			}
		} else {
			fmt.Printf("No hotkey requests were found in donation from %s\n", donation.Username)
		}
}

func HandleRequest(identifier string, args []string, money float64) (remains float64) {
	if command, ok := config.CommandList.Commands[identifier]; ok {
		switch command.Type {
		case "Hotkey":
			if len(args) != 1 {
				fmt.Printf("Should have seconds in args for hotkey command '%s', skipping...\n", identifier)
				return money
			}
			hotkeyId := command.Args.Id
			seconds, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid seconds value for hotkey command '%s', skipping...\n", identifier)
				return money
			}
			if seconds < DefaultPlanner.minimumDuration {
				fmt.Printf("Hotkey command '%s' was requested for %d seconds but it less than minimum duration %d seconds (can be changed in config), skipping...\n", identifier, seconds, DefaultPlanner.minimumDuration)
				return money
			}
			if float64(command.Price) * float64(seconds) > money {
				fmt.Printf("Not enough money for hotkey command '%s': need %.2f RUB but only %.2f RUB available, skipping...\n", identifier, float64(command.Price) * float64(seconds), money)
				return money
			}
			DefaultPlanner.AddHotkeyToSchedule(identifier, hotkeyId, seconds)
			return money - (float64(command.Price) * float64(seconds))
		case "Tint":
			if len(args) != 1 {
				fmt.Printf("Should have seconds in args for Tint command '%s', skipping...\n", identifier)
				return money
			}
			seconds, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid seconds value for Tint command '%s', skipping...\n", identifier)
				return money
			}
			if seconds < DefaultPlanner.minimumDuration {
				fmt.Printf("Tint command '%s' was requested for %d seconds but it less than minimum duration %d seconds (can be changed in config), skipping...\n", identifier, seconds, DefaultPlanner.minimumDuration)
				return money
			}
			if float64(command.Price) * float64(seconds) > money {
				fmt.Printf("Not enough money for Tint command '%s': need %.2f RUB but only %.2f RUB available, skipping...\n", identifier, float64(command.Price) * float64(seconds), money)
				return money
			}
			for colorName, color := range command.Args.Colors {
				go func() {
					err := DefaultPlanner.AddTintToSchedule(colorName, color.R, color.G, color.B, seconds, color.ListOfMeshes)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
						return
					}
				}()
			}
			return money - (float64(command.Price) * float64(seconds))
		default:
			fmt.Printf("Unknown command type '%s' for command '%s', skipping...\n", command.Type, identifier)
			return money
		}
	}
	fmt.Printf("Unknown command identifier '%s', skipping...\n", identifier)
	return money
}