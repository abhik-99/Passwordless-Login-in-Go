package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetENV(name string) string {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("ERROR while loading the env file")
		log.Fatal(err)
	}
	return os.Getenv(name)
}
