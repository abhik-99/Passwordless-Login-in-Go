package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var validate = validator.New()

type GenericJsonResponseDTO struct {
	Message string `json:"message"`
}

// DecodeJSONRequest decodes the JSON request body into the provided interface and validates it.
func DecodeJSONRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}

	// Validate the decoded struct
	return validate.Struct(v)
}

// EncodeJSONResponse encodes the provided interface as JSON and writes it to the response writer.
func EncodeJSONResponse(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

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
