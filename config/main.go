package config

import (
	"os"
)

type Config struct {
	DBPath string
}

func Get(getenv func(string) string) (Config, error) {
	if getenv == nil {
		getenv = os.Getenv
	}

	return Config{
		DBPath: getenv("DB_PATH"),
	}, nil
}
