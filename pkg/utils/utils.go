package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
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

func GenerateJWT(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		// Also fixed dates can be used for the NumericDate
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		Issuer:    GetENV("JWT_ISSUER"),
		ID:        id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetENV("JWT_SECRET"))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(GetENV("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, err
	}

}

func OTPGenerator() (string, error) {
	const otpChars = "0123456789"
	buffer := make([]byte, 8)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < 8; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func SendOTPMail(receipientEmail string, otp string) error {
	from := mail.NewEmail(GetENV("SENDER_NAME"), GetENV("SENDER_EMAIL")) // Change to your verified sender
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail(receipientEmail, receipientEmail)
	plainTextContent := fmt.Sprintf("Your OTP is %s", otp)
	htmlContent := fmt.Sprintf("Your OTP is <strong>%s</string>", otp)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	if _, err := client.Send(message); err != nil {
		log.Println("[ERROR]", err)
		return err
	}
	return nil

}

func SendOTPSms(receipientNo string, otp string) error {
	client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(fmt.Sprintf("Your OTP for Login is %s.", otp))
	params.SetFrom(GetENV("TWILLIO_PHONE_NO"))
	params.SetTo(receipientNo)

	_, err := client.Api.CreateMessage(params)
	return err
}
