package planner

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/makmanu/client_for_donatex/plugin"
)

type Planner struct {
	minimumDuration int
	plannerCh       chan string
	ticker          *time.Ticker
	schedule        map[string]time.Time
	oldSchedule     map[string]time.Duration
}

var DefaultPlanner *Planner

func NewPlanner(minimumDuration int) *Planner {
	plannerCh := make(chan string)
	DefaultPlanner = &Planner{
		minimumDuration: minimumDuration,
		plannerCh:       plannerCh,
		ticker:          time.NewTicker(100 * time.Millisecond),
		schedule:        make(map[string]time.Time),
		oldSchedule:     make(map[string]time.Duration),
	}
	return DefaultPlanner
}

func (p *Planner) Start() {
	signal := ""
	changed := false
	for {
		select {
		case <-p.ticker.C:
			now := time.Now()
			for key, t := range p.schedule {
				if now.After(t) {
					fmt.Printf("Executing hotkey %s\n", key)
					p.RemoveFromSchedule(key)
					continue
				}
				if oldT, ok := p.oldSchedule[key]; !ok || oldT.Round(time.Second).Seconds() != t.Sub(now).Round(time.Second).Seconds() {
					p.oldSchedule[key] = t.Sub(now)
					changed = true
				}
			}
			if changed {
				fmt.Printf("Current schedule:\n")
				for key, t := range p.schedule {
					fmt.Printf("  %s: %s\n", key, t.Sub(now).Round(time.Second))
				}
				changed = false
			}

		case signal = <-p.plannerCh:
			signalParts := strings.Fields(signal)
			switch signalParts[0] {
			case "add":
				if len(signalParts) < 3 {
					fmt.Println("Invalid add signal. Usage: add <key> <time_in_seconds>")
					return
				}
				key := signalParts[1]
				time_in_seconds, err := strconv.Atoi(signalParts[2])
				if err != nil {
					fmt.Println("Invalid time value. Please enter a valid integer.")
					return
				}
				err = p.AddToSchedule(key, time_in_seconds)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			default:
				fmt.Printf("Unknown signal: %s\n", signalParts[0])
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

func (p *Planner) RemoveFromSchedule(key string) {
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(key)
	if err != nil {
		return
	}
	delete(p.schedule, hotkey.Name)
	delete(p.oldSchedule, hotkey.Name)
	plugin.ExecuteHotkey(hotkey.Name)
}
