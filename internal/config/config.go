package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env:"ENV" env-default:"local"`
	HTTP HTTP   `yaml:"http"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return mustLoadFromEnv()
	}

	return mustLoadByPath(path)
}

func mustLoadByPath(configPath string) *Config {
	cfg, err := loadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadByPath(configPath string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func mustLoadFromEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Env empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	if flag.Lookup("config") == nil {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
