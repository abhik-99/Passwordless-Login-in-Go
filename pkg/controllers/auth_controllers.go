package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/abhik-99/passwordless-login/pkg/utils"
	"github.com/gorilla/mux"
)

func Signup(w http.ResponseWriter, r *http.Request) {}

func OTPViaEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if email, err := mail.ParseAddress(params["emailId"]); err != nil {
		log.Panic("Not a Valid Email")
	} else {
		fmt.Printf("Email is %s.\n", email)
	}
}
func OTPViaPhone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	phoneNo := params["phoneNo"]
	if !utils.IsValidPhoneNumber(phoneNo) {
		log.Panic("Not a Valid Phone")
	} else {
		fmt.Printf("Phone Number is %s.\n", phoneNo)
	}
}

func LoginViaEmail(w http.ResponseWriter, r *http.Request) {}
func LoginViaPhone(w http.ResponseWriter, r *http.Request) {}
