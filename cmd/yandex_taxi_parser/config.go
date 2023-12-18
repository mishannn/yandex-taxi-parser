package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Proxies struct {
		Type string `yaml:"type"`
		URL  string `yaml:"url"`
	} `yaml:"proxies"`
	Cookies struct {
		URL string `yaml:"url"`
	} `yaml:"cookies"`
	Points struct {
		Home Coordinates `yaml:"home"`
		Work Coordinates `yaml:"work"`
	} `yaml:"points"`
	Database struct {
		Address  string `yaml:"address"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func readConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, fmt.Errorf("can't parse config file: %w", err)
	}

	return config, nil
}
