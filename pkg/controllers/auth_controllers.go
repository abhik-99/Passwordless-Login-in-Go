package controllers

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/abhik-99/passwordless-login/pkg/data"
	"github.com/abhik-99/passwordless-login/pkg/utils"
	"github.com/gorilla/mux"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var newUser data.CreateUserDTO
	err := utils.DecodeJSONRequest(r, &newUser)
	if err != nil {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if _, err := data.CreateNewUser(newUser); err != nil {
		utils.EncodeJSONResponse(w, http.StatusInternalServerError, utils.GenericJsonResponseDTO{Message: "Error Occurred while Creating user"})
		return
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, utils.GenericJsonResponseDTO{Message: "User Created Successfully"})
	}

}

func OTPViaEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if email, err := mail.ParseAddress(params["emailId"]); err != nil {

		utils.EncodeJSONResponse(w, http.StatusBadRequest, err)
		return
	} else {
		if result, userID, lookUpErr := data.UserLookupViaEmail(email.String()); lookUpErr != nil {
			http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
			fmt.Println("[ERROR]", lookUpErr)
			return
		} else if !result {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		} else {
			otp, err := utils.OTPGenerator()
			if err != nil {
				http.Error(w, "Error Occurred while generating user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR]", err)
				return
			}
			a := data.Auth{UserId: userID, Otp: otp}
			if err := a.SetOTPForUser(); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR]", err)
				return
			}
			if err := utils.SendOTPMail(email.String(), otp); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR]", err)
				return
			}

			utils.EncodeJSONResponse(w, http.StatusOK, map[string]string{"message": "User OTP sent"})

		}
	}
}
func OTPViaPhone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	phoneNo := params["phoneNo"]
	if !utils.IsValidPhoneNumber(phoneNo) {
		http.Error(w, "Invalid Phone number", http.StatusBadRequest)
		return
	} else {
		if result, userID, lookUpErr := data.UserLookupViaPhone(phoneNo); lookUpErr != nil {
			http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
			fmt.Println("[ERROR]", lookUpErr)
			return
		} else if !result {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		} else {
			a := data.Auth{UserId: userID, Otp: "123456"}
			if err := a.SetOTPForUser(); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR]", err)
				return
			}
			utils.EncodeJSONResponse(w, http.StatusOK, map[string]string{"message": "User OTP sent"})

		}
	}
}

func LoginViaEmail(w http.ResponseWriter, r *http.Request) {
	var dto data.LoginWithEmailDTO
	err := utils.DecodeJSONRequest(r, dto)
	if err != nil {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, err)
		return
	}
	result, userID, lookUpErr := data.UserLookupViaEmail(dto.Email)

	if !result {
		utils.EncodeJSONResponse(w, http.StatusNotFound, "User Does not Exist")
		return
	}
	if lookUpErr != nil {
		http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
		fmt.Println("[ERROR]", lookUpErr)
		return
	}
	a := data.Auth{UserId: userID, Otp: dto.Otp}
	if result, err := a.CheckOTP(); err != nil {
		http.Error(w, "Error Ocurred while verifying OTP", http.StatusInternalServerError)
		fmt.Println("[ERROR]", err)
		return
	} else if !result {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, "OTP Does not Match")
		return
	}

	if result, err := utils.GenerateJWT(userID); err != nil {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, utils.GenericJsonResponseDTO{Message: "Error Occured while access token generation "})
		return
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, data.AccessTokenDTO{AccessToken: result})
	}

}
func LoginViaPhone(w http.ResponseWriter, r *http.Request) {}
