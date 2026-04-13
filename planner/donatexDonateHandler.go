package planner

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/makmanu/client_for_donatex/plugin"
	"github.com/makmanu/client_for_donatex/structs"
)

func GetRequestedHotkeysFromMessage(message string) (map[string]int, error) {
	result := make(map[string]int)
	
	re := regexp.MustCompile(`&([^#]+)#([0-9]+)`)
	matches := re.FindAllStringSubmatch(message, -1)

	for _, m := range matches {
		identifier := m[1]
		number := m[2]

		// normalize number
		num, _ := strconv.Atoi(number)

		// normalize identifier ONLY if it's purely numeric
		if idNum, err := strconv.Atoi(identifier); err == nil {
			identifier = strconv.Itoa(idNum)
		}

		result[identifier] = num
	}

	return result, nil
}

func DonationIsBigEnough(reqHotkeys map[string]int, donationAmount float64) (bool, error) {
	var sum float64
	for identifier, seconds := range reqHotkeys {
		if seconds < DefaultPlanner.minimumDuration {
			delete(reqHotkeys, identifier)
			continue
		}
		hotkey, err := plugin.FindHotkeyInfoByIdentifier(identifier)
		if err != nil {
			delete(reqHotkeys, identifier)
			continue
		}
		sum += hotkey.Coefficient * float64(seconds)
	}
	return donationAmount >= sum, nil
}

func AddHotkeysToSchedule(reqHotkeys map[string]int) {
	for identifier, seconds := range reqHotkeys {
		DefaultPlanner.plannerCh <- fmt.Sprintf("add %s %d", identifier, seconds)
		log.Printf("Added hotkey %s for %d seconds to schedule\n", identifier, seconds)
	}
}

func HandleDonation(donation structs.Donation) {
	reqHotkeys, err := GetRequestedHotkeysFromMessage(donation.Message)
		if err != nil {
			fmt.Printf("Error occurred while parsing hotkeys: %v\n", err)
		}

		if len(reqHotkeys) > 0 {
			enough, err := DonationIsBigEnough(reqHotkeys, donation.AmountInRub)
			if err != nil {
				fmt.Printf("Error occurred while checking donation amount: %v\n", err)
			}
			if !enough {
				fmt.Printf(" Donation from %s is not enough for requested hotkeys\n", donation.Username)
			} else {
				AddHotkeysToSchedule(reqHotkeys)
			}
		} else {
			fmt.Printf(" No hotkeys requested in donation from %s\n", donation.Username)
		}
}