package main

import (
	"log"
	"os"
)

var (
	CLIENT_ID     string
	CLIENT_SECRET string
	ACCESS_TOKEN  string
	REFRESH_TOKEN string
)

func getEnv(env string) string {
	env, exists := os.LookupEnv(env)
	if !exists {
		log.Fatalf("Environment variable %v doesn't exist", env)
	}

	return env
}

func getUncheckedEnv(env string) string {
	env, _ = os.LookupEnv(env)

	return env

}

func init() {
	CLIENT_ID = getEnv("CLIENT_ID")
	CLIENT_SECRET = getEnv("CLIENT_SECRET")
	ACCESS_TOKEN = getUncheckedEnv("CLIENT_SECRET")
	REFRESH_TOKEN = getUncheckedEnv("CLIENT_SECRET")
}
