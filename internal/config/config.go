package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(current_user_name string) error {
	c.CurrentUserName = current_user_name

	err := write(*c)
	if err != nil {
	    return err
	}

	return nil
}

func Read() (Config, error) {
	var config Config

	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening file: %w", err)
	}
	defer jsonFile.Close()

	err = json.NewDecoder(jsonFile).Decode(&config)

	if err != nil {
		return Config{}, fmt.Errorf("error deconding json: %w", err)
	}

	return config, nil
}

func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
	    return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(cfg)
	if err != nil {
	    return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homepath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home path: %w", err)
	}

	return homepath + "/" + configFileName, nil
}
