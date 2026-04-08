package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type VTubeStudioConfig struct {
	URL        string `yaml:"url"`
	Port       int    `yaml:"port"`
	PluginName string `yaml:"pluginName"`
}

type Config struct {
	URL         string            `yaml:"url"`
	Port        int               `yaml:"port"`
	CertFile    string            `yaml:"certFile"`
	KeyFile     string            `yaml:"keyFile"`
	LogFile     string            `yaml:"logFile"`
	VTubeStudio VTubeStudioConfig `yaml:"VTubeStudio"`
}

type Secret struct {
	Token         string `yaml:"token"`
	WebhookSecret string `yaml:"webhookSecret"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadSecret(filename string) (*Secret, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var secret Secret
	err = yaml.Unmarshal(data, &secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}
