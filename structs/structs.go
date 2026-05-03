package structs

import "time"

type VTubeStudioConfig struct {
	URL        string `yaml:"url"`
	Port       int    `yaml:"port"`
	PluginName string `yaml:"pluginName"`
}

type Config struct {
	URL                        string            `yaml:"url"`
	VTubeStudio                VTubeStudioConfig `yaml:"VTubeStudio"`
	MinimumDuration            int               `yaml:"minimumDuration"`
	MaximumHotkeysPerDonation  int               `yaml:"maximumHotkeysPerDonation"`
	DefaultCoefficient         float64           `yaml:"defaultHotkeyCoefficient"`
}

type Secret struct {
	Token         string `yaml:"token"`
}

type Hotkey struct {
	ID          int     `yaml:"id"`
	Name        string  `yaml:"name"`
	HotkeyID    string  `yaml:"hotkeyID"`
	Coefficient float64 `yaml:"coefficient"`
}

type Donation struct {
	ID                  string  `json:"id"`
	Username            string  `json:"username"`
	Message             string  `json:"message"`
	Currency            string  `json:"currency"`
	Amount              float64 `json:"amount"`
	AmountInRub         float64 `json:"amountInRub"`
	Timestamp           string  `json:"timestamp"`
	WithAiResponse      bool    `json:"withAiResponse"`
	AiResponse          string  `json:"aiResponse"`
	IsTest              bool    `json:"isTest"`
	IsPotentiallyUnsafe bool    `json:"isPotentiallyUnsafe"`
	WasShown            bool    `json:"wasShown"`
	IsFeePaidByUser     bool    `json:"isFeePaidByUser"`
	VoiceFilePath       string  `json:"voiceFilePath"`
	PaidVoice           string  `json:"paidVoice"`
	MusicLink           string  `json:"musicLink"`
}

type PlannerSignal struct {
	Command  string
	Item     ScheduleItem
}

type ScheduleItem struct {
	Name              string
	ArgsforExecute    []any
	ArgsforStop       []any
	FunctionToExecute func(...any) error
	FunctionToStop	  func(...any) error
	DurationinSeconds int
	ExpiresAt         time.Time
}

type ArtMeshInfo struct {
	ModelLoaded          bool
	NumberOfArtMeshNames int
	NumberOfArtMeshTags  int
	ArtMeshNames         []string
	ArtMeshTags          []string
}

type Command struct {
	Name  string      `yaml:"Name"`
	Type  string      `yaml:"Type"`
	Price float64     `yaml:"Price"`
	Args  CommandArgs `yaml:"Args"`
}

type CommandArgs struct {
	Tags []string `yaml:"Tags,omitempty"`
	Id   string   `yaml:"Id,omitempty"`
	Colors map[string]Color `yaml:",inline"`
}

type Color struct {
	ListOfMeshes []string `yaml:"ListOfMeshes"`
	R            int      `yaml:"R"`
	G            int      `yaml:"G"`
	B            int      `yaml:"B"`
}

type CommandList struct {
	Commands map[string]Command `yaml:"Commands"`
}