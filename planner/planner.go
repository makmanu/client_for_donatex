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
	schedule                  map[string]time.Time
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
		schedule:                  make(map[string]time.Time),
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
				if now.After(t) {
					err := p.RemoveFromSchedule(key)
					if err != nil {
						log.Printf("Error removing hotkey from schedule: %v\n", err)
					}
					continue
				}
				if oldT, ok := p.oldSchedule[key]; !ok || oldT.Round(time.Second).Seconds() != t.Sub(now).Round(time.Second).Seconds() {
					p.oldSchedule[key] = t.Sub(now)
					changed = true
				}
			}
			if changed {
				p.scheduleScreen.Truncate(0)
				p.scheduleScreen.Seek(0, 0)
				p.scheduleScreen.WriteString("Current schedule:\n")
				for key, t := range p.schedule {
					if secondsToDisplay := t.Sub(now).Round(time.Second); secondsToDisplay > 0 {
						p.scheduleScreen.WriteString(fmt.Sprintf("  %s: %s\n", key, secondsToDisplay))
					}
				}
				changed = false
			}

		case signal := <-p.plannerCh:
			switch signal.Command {
			case "add":
				err := p.AddToSchedule(signal.Hotkey, signal.Seconds)
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

func (p *Planner) AddToSchedule(key string, time_in_seconds int) error {
	t := time.Duration(time_in_seconds) * time.Second
	if t < (time.Duration(p.minimumDuration) * time.Second) {
		return fmt.Errorf("Duration must be at least %d seconds", p.minimumDuration)
	}
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(key)
	if err != nil {
		return fmt.Errorf("Error finding hotkey: %w", err)
	}
	if _, ok := p.schedule[hotkey.Name]; !ok || time.Now().After(p.schedule[hotkey.Name]) {
		p.schedule[hotkey.Name] = time.Now().Add(t)
		plugin.ExecuteHotkey(hotkey.Name)
	} else {
		p.schedule[hotkey.Name] = p.schedule[hotkey.Name].Add(t)
	}
	return nil
}

func (p *Planner) RemoveFromSchedule(key string) error {
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(key)
	if err != nil {
		return fmt.Errorf("Error finding hotkey: %w", err)
	}
	for hotkeyName := range p.schedule {
		if hotkeyName == hotkey.Name {
			delete(p.schedule, hotkeyName)
			delete(p.oldSchedule, hotkeyName)
			err = plugin.ExecuteHotkey(hotkeyName)
			if err != nil {
				return fmt.Errorf("Error executing hotkey: %w", err)
			}
			fmt.Printf("Removed hotkey %s from schedule\n", hotkeyName)
			return nil
		}
	}
	fmt.Printf("Hotkey '%s' not found in schedule", key)
	return nil
}
