package configs

import (
	"log/slog"
	"os"
)

type ConfigList struct {
	Env                  string
	APICorsAllowsOrigins []string
}

var Configs ConfigList

func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

func LoadEnv() error {
	Configs = ConfigList{
		Env:                  GetEnvDefault("APP_ENV", "development"),
		APICorsAllowsOrigins: []string{"http://localhost:3001"},
	}
	return nil
}

func GetEnvDefault(key, defVal string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defVal
	}
	return val
}

func init() {
	if err := LoadEnv(); err != nil {
		slog.Error("Failed to load env", "err", err)
		panic(err)
	}
	slog.Debug("Success init env")
}
