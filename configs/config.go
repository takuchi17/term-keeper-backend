package configs

import (
	"log/slog"
	"os"
	"strconv"
)

type ConfigList struct {
	Env                  string
	DBUser               string
	DBHost               string
	DBPort               int
	DBName               string
	DBPassword           string
	APICorsAllowsOrigins []string
	JWTSecret            string
}

var Config ConfigList

func (c *ConfigList) IsDevelopment() bool {
	return c.Env == "development"
}

func LoadEnv() error {
	DBPort, err := strconv.Atoi(getEnvDefault("MYSQL_PORT", "3307"))
	if err != nil {
		return nil
	}

	Config = ConfigList{
		Env:                  getEnvDefault("APP_ENV", "development"),
		DBUser:               getEnvDefault("DB_USER", "user"),
		DBHost:               getEnvDefault("DB_HOST", "localhost"),
		DBPort:               DBPort,
		DBName:               getEnvDefault("DB_NAME", "term_keeper_db"),
		DBPassword:           getEnvDefault("DB_PASSWORD", "password"),
		APICorsAllowsOrigins: []string{"http://localhost:3001"},
		JWTSecret:            getEnvDefault("JWT_SECRET", "default-secret"),
	}
	return nil
}

func getEnvDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
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
