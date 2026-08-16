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
	minimumDuration            int
	maximumRequestsPerDonation int
	plannerCh                  chan structs.PlannerSignal
	ticker                     *time.Ticker
	schedule                   map[string]*structs.ScheduleItem
	oldSchedule                map[string]time.Duration
	scheduleScreen             *os.File
}

var DefaultPlanner *Planner

func NewPlanner(minimumDuration, maximumRequestsPerDonation int) *Planner {
	plannerCh := make(chan structs.PlannerSignal)
	Screen, err := os.OpenFile("scheduleScreen.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening schedule screen file: %v\n", err)
	} else {
		Screen.WriteString("Current schedule:\n")
	}
	DefaultPlanner = &Planner{
		maximumRequestsPerDonation: maximumRequestsPerDonation,
		minimumDuration:           	minimumDuration,
		plannerCh:                 	plannerCh,
		ticker:                    	time.NewTicker(100 * time.Millisecond),
		schedule:                  	make(map[string]*structs.ScheduleItem),
		oldSchedule:               	make(map[string]time.Duration),
		scheduleScreen:            	Screen,
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
			item := p.schedule[itemName]
			if item != nil && item.FunctionToStop != nil {
				go func(it *structs.ScheduleItem) {
					err := it.FunctionToStop(it.ArgsforStop...)
					if err != nil {
						log.Printf("Error executing function: %v", err)
					}
				}(item)
			}
			delete(p.schedule, itemName)
			delete(p.oldSchedule, itemName)
			fmt.Printf("Removed item %s from schedule\n", itemName)
			return nil
		}
	}
	fmt.Printf("Item '%s' not found in schedule", itemNameToRemove)
	return nil
}

func (p *Planner) AddHotkeyToSchedule(name string, hotkeyId string, seconds int) error {
	item := &structs.ScheduleItem{
		Name:           name,
		ArgsforExecute: []any{hotkeyId},
		ArgsforStop:    []any{hotkeyId},
		FunctionToExecute: func(args ...any) error {
			return plugin.ExecuteHotkey(args[0].(string))
		},
		FunctionToStop: func(args ...any) error {
			return plugin.ExecuteHotkey(args[0].(string))
		},
		DurationinSeconds: seconds,
		ExpiresAt: time.Now(),
	}
	err := p.AddItemToSchedule(item)
	if err != nil {
		return fmt.Errorf("Error adding hotkey to schedule: %w", err)
	}
	log.Printf("Added hotkey %s for %d seconds to schedule\n", hotkeyId, seconds)
	return nil
}

func (p *Planner) RemoveHotkeyFromSchedule(name string) error {
	return p.RemoveItemFromSchedule(name)
}

func (p *Planner) AddTintToSchedule(name string, r, g, b, a int, seconds int, meshes []string) error {
	item := &structs.ScheduleItem{
		Name: name,
		ArgsforExecute: []any{r, g, b, a, 1, meshes},
		ArgsforStop:    []any{r, g, b, a, 0, meshes},
		FunctionToExecute: func(args ...any) error {
			return plugin.TintMeshesFadeIn(255, 255, 255, 255, args[0].(int), args[1].(int), args[2].(int), args[3].(int), args[4].(int), args[5].([]string))
		},
		FunctionToStop: func(args ...any) error {
			return plugin.TintMeshesFadeIn(args[0].(int), args[1].(int), args[2].(int), args[3].(int), 255, 255, 255, 255, args[4].(int), args[5].([]string))
		},
		DurationinSeconds: seconds,
		ExpiresAt: time.Now(),
	}
	err := p.AddItemToSchedule(item)
	if err != nil {
		return fmt.Errorf("Error adding tint to schedule: %w", err)
	}
	log.Printf("Added tint RGBA(%d, %d, %d, %d) for %d seconds to schedule\n", r, g, b, 255, seconds)
	return nil
}