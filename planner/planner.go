package planner

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/makmanu/client_for_donatex/plugin"
	"github.com/makmanu/client_for_donatex/structs"
)

type Planner struct {
	minimumDuration           int
	maximumHotkeysPerDonation int
	plannerCh                 chan structs.PlannerSignal
	ticker                    *time.Ticker
	schedule                  map[string]*structs.ScheduleItem
	oldSchedule               map[string]time.Duration
	scheduleScreen            *os.File
}

var DefaultPlanner *Planner

func NewPlanner(minimumDuration, maximumHotkeysPerDonation int) *Planner {
	plannerCh := make(chan structs.PlannerSignal)
	Screen, err := os.OpenFile("scheduleScreen.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening schedule screen file: %v\n", err)
	} else {
		Screen.WriteString("Current schedule:\n")
	}
	DefaultPlanner = &Planner{
		maximumHotkeysPerDonation: maximumHotkeysPerDonation,
		minimumDuration:           minimumDuration,
		plannerCh:                 plannerCh,
		ticker:                    time.NewTicker(100 * time.Millisecond),
		schedule:                  make(map[string]*structs.ScheduleItem),
		oldSchedule:               make(map[string]time.Duration),
		scheduleScreen:            Screen,
	}
	return DefaultPlanner
}

func (p *Planner) Start() {
	changed := false
	for {
		select {
		case <-p.ticker.C:
			now := time.Now()
			for key, t := range p.schedule {
				if now.After(t.ExpiresAt) {
					err := p.RemoveItemFromSchedule(key)
					if err != nil {
						log.Printf("Error removing hotkey from schedule: %v\n", err)
					}
					continue
				}
				if oldT, ok := p.oldSchedule[key]; !ok || oldT.Round(time.Second).Seconds() != t.ExpiresAt.Sub(now).Round(time.Second).Seconds() {
					p.oldSchedule[key] = t.ExpiresAt.Sub(now)
					changed = true
				}
			}
			if changed {
				p.scheduleScreen.Truncate(0)
				p.scheduleScreen.Seek(0, 0)
				p.scheduleScreen.WriteString("Current schedule:\n")
				for key, t := range p.schedule {
					if secondsToDisplay := t.ExpiresAt.Sub(now).Round(time.Second); secondsToDisplay > 0 {
						p.scheduleScreen.WriteString(fmt.Sprintf("  %s: %s\n", key, secondsToDisplay))
					}
				}
				changed = false
			}

		case signal := <-p.plannerCh:
			switch signal.Command {
			case "add":
				err := p.AddItemToSchedule(&signal.Item)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			default:
				fmt.Printf("Unknown signal: %s\n", signal.Command)
			}
		}
	}
}

func (p *Planner) AddItemToSchedule(item *structs.ScheduleItem) error {
	t := time.Duration(item.DurationinSeconds) * time.Second

	if t < (time.Duration(p.minimumDuration) * time.Second) {
		return fmt.Errorf("Duration must be at least %d seconds", p.minimumDuration)
	}

	now := time.Now()
	if _, ok := p.schedule[item.Name]; !ok  {
		p.schedule[item.Name] = item
		p.schedule[item.Name].ExpiresAt = now.Add(t)
		err := item.FunctionToExecute(item.ArgsforExecute...)
		if err != nil {
			return fmt.Errorf("Error executing function: %w", err)
		}
	} else {
		p.schedule[item.Name].ExpiresAt = p.schedule[item.Name].ExpiresAt.Add(t)
	}
	return nil
}

func (p *Planner) RemoveItemFromSchedule(itemNameToRemove string) error {

	for itemName := range p.schedule {
		if itemName == itemNameToRemove {
			err := p.schedule[itemName].FunctionToStop(p.schedule[itemName].ArgsforStop...)
			if err != nil {
				log.Printf("Error executing function: %v", err)
			}
			delete(p.schedule, itemName)
			delete(p.oldSchedule, itemName)
			fmt.Printf("Removed item %s from schedule\n", itemName)
			return err
		}
	}
	fmt.Printf("Item '%s' not found in schedule", itemNameToRemove)
	return nil
}

func (p *Planner) AddHotkeyToSchedule(identifier string, seconds int) error {
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(identifier)
	if err != nil {
		return fmt.Errorf("Couldn't find hotkey with identifier '%s': %v", identifier, err)
	}
	item := &structs.ScheduleItem{
		Name:           hotkey.Name,
		ArgsforExecute: []any{hotkey},
		ArgsforStop:    []any{hotkey},
		FunctionToExecute: func(args ...any) error {
			return plugin.ExecuteHotkey(args[0].(structs.Hotkey))
		},
		FunctionToStop: func(args ...any) error {
			return plugin.ExecuteHotkey(args[0].(structs.Hotkey))
		},
		DurationinSeconds: seconds,
		ExpiresAt: time.Now(),
	}
	err = p.AddItemToSchedule(item)
	if err != nil {
		return fmt.Errorf("Error adding hotkey to schedule: %w", err)
	}
	log.Printf("Added hotkey %s for %d seconds to schedule\n", identifier, seconds)
	return nil
}

func (p *Planner) RemoveHotkeyFromSchedule(hotkey structs.Hotkey) error {
	return p.RemoveItemFromSchedule(hotkey.Name)
}