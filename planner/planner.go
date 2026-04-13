package planner

import (
	"fmt"
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
}

var DefaultPlanner *Planner

func NewPlanner(minimumDuration, maximumHotkeysPerDonation int) *Planner {
	plannerCh := make(chan structs.PlannerSignal)
	DefaultPlanner = &Planner{
		maximumHotkeysPerDonation: maximumHotkeysPerDonation,
		minimumDuration:           minimumDuration,
		plannerCh:                 plannerCh,
		ticker:                    time.NewTicker(100 * time.Millisecond),
		schedule:                  make(map[string]time.Time),
		oldSchedule:               make(map[string]time.Duration),
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

func (p *Planner) RemoveFromSchedule(key string) {
	hotkey, err := plugin.FindHotkeyInfoByIdentifier(key)
	if err != nil {
		return
	}
	delete(p.schedule, hotkey.Name)
	delete(p.oldSchedule, hotkey.Name)
	plugin.ExecuteHotkey(hotkey.Name)
}
