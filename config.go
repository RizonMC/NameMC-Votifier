package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	NameMCAddress string `json:"serverAddress"`
	Votifier      struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
		Token   string `json:"token"`
	} `json:"votifier"`
}

func ReadConfig(name string) (*Config, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
