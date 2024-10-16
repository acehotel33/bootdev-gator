package config

import (
	"encoding/json"
	"os"
)

const cfgFile = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func InitializeConfig() (*Config, error) {
	cfg, err := Read()
	if err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUsername = user

	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(cfgPath, cfgJSON, 0644); err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgContents, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(cfgContents, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := homePath + "/" + cfgFile

	return path, nil
}
