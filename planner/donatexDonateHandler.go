package planner

import (
	"fmt"
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
			fmt.Printf("Couldn't find hotkey with identifier '%s': %v\n", identifier, err)
			delete(reqHotkeys, identifier)
			continue
		}
		fmt.Printf("Hotkey %s with coefficient %.2f requested for %d seconds(%.2f roubles)\n", hotkey.Name, hotkey.Coefficient, seconds, hotkey.Coefficient * float64(seconds))
		sum += hotkey.Coefficient * float64(seconds)
	}
	fmt.Printf("Requested amount %.2f, donated amount %.2f\n", sum, donationAmount)
	return donationAmount >= sum, nil
}

func HandleDonation(donation structs.Donation) {
	reqHotkeys, err := GetRequestedHotkeysFromMessage(donation.Message)
		if err != nil {
			fmt.Printf("Error occurred while parsing hotkeys: %v\n", err)
		}
		if len(reqHotkeys) > DefaultPlanner.maximumHotkeysPerDonation {
			fmt.Printf("Too many hotkeys requested: %d. Maximum is %d(can be changed in config), skipping...\n", len(reqHotkeys), DefaultPlanner.maximumHotkeysPerDonation)
		}
		if len(reqHotkeys) > 0 {
			donationAmount := 0.0
			if donation.Currency == "RUB" {
				fmt.Println(donation.Amount)
				donationAmount = donation.Amount
			} else {
				fmt.Println("false")
				donationAmount = donation.AmountInRub
			}
			fmt.Printf("? Donation amount: %.2f RUB\n", donationAmount)
			for identifier, seconds := range reqHotkeys {
				if seconds < DefaultPlanner.minimumDuration {
					fmt.Printf("- Hotkey '%s' was requested for %d seconds but it less than minimum duration %d seconds (can be changed in config), skipping...\n", identifier, seconds, DefaultPlanner.minimumDuration)
					continue
				}
				hotkey, err := plugin.FindHotkeyInfoByIdentifier(identifier)
				if err != nil {
					fmt.Printf("- Couldn't find hotkey with identifier '%s', skipping...\n", identifier)
					continue
				}
				needToSpent := hotkey.Coefficient * float64(seconds)
				if donationAmount < needToSpent {
					fmt.Printf("- Not enough money for hotkey '%s'. Need %.2f RUB but %.2f RUB remain, skipping...\n", hotkey.Name, needToSpent, donationAmount)
					continue
				}
				donationAmount -= needToSpent
				fmt.Printf("✓ Hotkey '%s' with coefficient %.2f requested for %d seconds(%.2f RUB)\n? Remain amount: %.2f RUB\n", hotkey.Name, hotkey.Coefficient, seconds, needToSpent, donationAmount)
				DefaultPlanner.AddHotkeyToSchedule(hotkey.Name, seconds)
			}
		} else {
			fmt.Printf("No hotkey requests were found in donation from %s\n", donation.Username)
		}
}