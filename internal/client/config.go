package client

import (
	"os"
	"time"
)

type Config struct {
	ServerAddr string
	Workers    int
	Timeout    time.Duration
}

func NewConfig() *Config {
	c := new(Config)
	c.ServerAddr = getEnvDefault("SERVER_ADDR", "127.0.0.1:9000")
	c.Workers = 4
	c.Timeout = time.Millisecond * time.Duration(1000)
	return c
}

func getEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}
