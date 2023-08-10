package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `json:"port"`
	DB   struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Name     string `json:"name"`
	} `json:"db"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Port: viper.GetString("port"),
		DB: struct {
			User     string `json:"user"`
			Password string `json:"password"`
			Host     string `json:"host"`
			Name     string `json:"name"`
		}{
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			Host:     viper.GetString("db.host"),
			Name:     viper.GetString("db.name"),
		},
	}

	return cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Port == "" {
		return errors.New("port must not be empty")
	}
	if cfg.DB.User == "" {
		return errors.New("db.user must not be empty")
	}
	if cfg.DB.Password == "" {
		return errors.New("db.password must not be empty")
	}
	if cfg.DB.Host == "" {
		return errors.New("db.host must not be empty")
	}
	if cfg.DB.Name == "" {
		return errors.New("db.name must not be empty")
	}
	return nil
}
