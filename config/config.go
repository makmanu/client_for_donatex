package config

import (
	"os"

	"github.com/makmanu/client_for_donatex/structs"
	"gopkg.in/yaml.v3"
)

var (
	config           structs.Config
	CommandList		 *structs.CommandList
)

func LoadConfig(filename string) (*structs.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadSecret(filename string) (*structs.Secret, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var secret structs.Secret
	err = yaml.Unmarshal(data, &secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func LoadCommandList(filename string) (*structs.CommandList, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &CommandList)
	if err != nil {
		return nil, err
	}
	return CommandList, nil
}