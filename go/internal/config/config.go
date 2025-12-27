package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env     string
	AppPort int

	MySQLHost     string
	MySQLPort     int
	MySQLDatabase string
	MySQLUser     string
	MySQLPassword string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() (Config, error) {
	var c Config
	c.Env = getEnv("APP_ENV", "development")
	c.AppPort = mustInt(getEnv("APP_PORT", "8080"))

	c.MySQLHost = getEnv("MYSQL_HOST", "mysql")
	c.MySQLPort = mustInt(getEnv("MYSQL_PORT", "3306"))
	c.MySQLDatabase = getEnv("MYSQL_DATABASE", "appdb")
	c.MySQLUser = getEnv("MYSQL_USER", "app")
	c.MySQLPassword = getEnv("MYSQL_PASSWORD", "app")

	c.RedisAddr = getEnv("REDIS_ADDR", "redis:6379")
	c.RedisPassword = getEnv("REDIS_PASSWORD", "")
	c.RedisDB = mustInt(getEnv("REDIS_DB", "0"))

	if c.MySQLHost == "" || c.MySQLDatabase == "" || c.MySQLUser == "" {
		return Config{}, fmt.Errorf("missing mysql config")
	}
	if c.RedisAddr == "" {
		return Config{}, fmt.Errorf("missing redis addr")
	}

	return c, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
