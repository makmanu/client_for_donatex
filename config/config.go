package config

import (
	"os"

	"github.com/makmanu/client_for_donatex/structs"
	"gopkg.in/yaml.v3"
)

var (
	Config           structs.Config
	CommandList		 *structs.CommandList
	Secret 			 *structs.Secret
)

func LoadConfig(filename string) (*structs.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return nil, err
	}

	return &Config, nil
}

func LoadSecret(filename string) (*structs.Secret, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &Secret)
	if err != nil {
		return nil, err
	}

	return Secret, nil
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