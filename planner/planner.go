package planner

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Planner struct {
	minimumDuration  float64
	plannerCh		 chan string
	ticker           *time.Ticker
	schedule         map[string]time.Time
}

func NewPlanner(minimumDuration float64) *Planner {
	plannerCh := make(chan string)
	return &Planner{
		minimumDuration:  minimumDuration,
		plannerCh:        plannerCh,
		ticker:           time.NewTicker(100 * time.Millisecond),
		schedule:         make(map[string]time.Time),
	}
}

func (p *Planner) Start() {
	signal := ""
	select {
	case <-p.ticker.C:
		now := time.Now()
		for key, t := range p.schedule {
			if now.After(t) {
				fmt.Printf("Executing hotkey %s\n", key)
				delete(p.schedule, key)
			}
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
			err = p.AddToShedule(key, time_in_seconds)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		default:
			fmt.Printf("Unknown signal: %s\n", signalParts[0])
		}
	}
}

func (p *Planner) AddToShedule(key string, time_in_seconds int) error{
	t := time.Duration(time_in_seconds) * time.Second
	if t < (time.Duration(p.minimumDuration) * time.Second) {
		return fmt.Errorf("Duration must be at least %f seconds", p.minimumDuration)
	}
	if _, ok := p.schedule[key]; !ok || time.Now().After(p.schedule[key]) {
		p.schedule[key] = time.Now().Add(t)
	} else {
		p.schedule[key] = p.schedule[key].Add(t)
	}
	return nil
}

func (p *Planner) RemoveFromSchedule(key string) {
	delete(p.schedule, key)

}