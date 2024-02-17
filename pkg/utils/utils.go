package utils

import (
	"log"
	"os"
	"regexp"

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

func IsValidPhoneNumber(phoneNo string) bool {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)

	return re.Find([]byte(phoneNo)) != nil
}
