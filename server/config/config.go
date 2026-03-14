package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Host    string
	Port    string
	RPCPath string
}

func (c config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func LoadConfig() config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return config{
		Host:    getEnv("SERVER_HOST"),
		Port:    getEnv("SERVER_PORT"),
		RPCPath: getEnv("SERVER_RPC_PATH"),
	}
}

func getEnv(envName string) string {
	env, exists := os.LookupEnv(envName)

	if !exists {
		log.Fatalf("Missing env variable %s", envName)
	}

	return env
}
