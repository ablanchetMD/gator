package config

import (
	"encoding/json"	
	"io"
	"os"

	"github.com/ablanchetmd/gator/internal/database"
)
const configFileName = ".gatorconfig.json"


func ReadConfig() (Config, error) {
	
	cfg := Config{}
	path, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}
	file, err := os.Open(path)
    if err != nil {
        return cfg, err
    }
    defer file.Close()

	// Read the entire file into memory
    data, err := io.ReadAll(file)
    if err != nil {
        return cfg, err
    }

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	err := writeConfig(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {

	//current, err := os.Getwd()
	current,err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := current + "/"+configFileName

	return path, nil
}

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Config *Config
	Db    *database.Queries
	Status string
}

