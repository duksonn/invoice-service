package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Load() (*Config, error) {
	cfg := &Config{}
	err := readConfig("dev", cfg)
	if err != nil {
		panic(fmt.Sprintf("config: error when read configuration from file | %s", err.Error()))
	}
	fmt.Println("config read successfully")

	return cfg, nil
}

func readConfig(env string, cfg *Config) error {
	err := unmarshalFromFile(fmt.Sprintf("cmd/config/env_%s.json", env), cfg)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		file, err := os.Open(fmt.Sprintf("env_%s.json", env))
		if err == nil {
			decoder := json.NewDecoder(file)
			err = decoder.Decode(cfg)
		}
		return err
	}
	return err
}

func unmarshalFromFile(filePath string, out interface{}) error {
	file, err := os.Open(filePath)

	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(out)
	}
	return err
}
