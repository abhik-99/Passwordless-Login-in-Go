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
		utils.EncodeJSONResponse(w, http.StatusBadRequest, struct {
			utils.GenericJsonResponseDTO
			Err string `json:"err"`
		}{
			GenericJsonResponseDTO: utils.GenericJsonResponseDTO{
				Message: "Invalid Request Body",
			},
			Err: err.Error(),
		})
		return
	}
	if _, err := data.CreateNewUser(newUser); err != nil {
		http.Error(w, "Error Occurred while Creating user", http.StatusInternalServerError)
		return
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, utils.GenericJsonResponseDTO{Message: "User Created Successfully"})
	}

}

func OTPViaEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if email, err := mail.ParseAddress(params["emailId"]); err != nil {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	} else {
		if result, userID, lookUpErr := data.UserLookupViaEmail(email.String()); lookUpErr != nil {
			http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
			fmt.Println("[ERROR] ", lookUpErr)
			return
		} else if !result {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		} else {
			otp, err := utils.OTPGenerator()
			if err != nil {
				http.Error(w, "Error Occurred while generating user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}
			a := data.Auth{UserId: userID, Otp: otp}
			if err := a.SetOTPForUser(); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}
			if err := utils.SendOTPMail(email.String(), otp); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}

			utils.EncodeJSONResponse(w, http.StatusOK, utils.GenericJsonResponseDTO{Message: "User OTP sent"})

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
			fmt.Println("[ERROR] ", lookUpErr)
			return
		} else if !result {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		} else {
			otp, err := utils.OTPGenerator()
			if err != nil {
				http.Error(w, "Error Occurred while generating user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}
			a := data.Auth{UserId: userID, Otp: otp}
			if err := a.SetOTPForUser(); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}
			if err := utils.SendOTPSms(phoneNo, otp); err != nil {
				http.Error(w, "Error Occurred while setting user OTP", http.StatusInternalServerError)
				fmt.Println("[ERROR] ", err)
				return
			}

			utils.EncodeJSONResponse(w, http.StatusOK, utils.GenericJsonResponseDTO{Message: "User OTP sent"})

		}
	}
}

func LoginViaEmail(w http.ResponseWriter, r *http.Request) {
	var dto data.LoginWithEmailDTO
	err := utils.DecodeJSONRequest(r, dto)
	if err != nil {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, struct {
			utils.GenericJsonResponseDTO
			Err string `json:"err"`
		}{
			GenericJsonResponseDTO: utils.GenericJsonResponseDTO{
				Message: "Invalid Request Body",
			},
			Err: err.Error(),
		})
		return
	}
	result, userID, lookUpErr := data.UserLookupViaEmail(dto.Email)

	if !result {
		http.Error(w, "User Does not Exist", http.StatusNotFound)
		return
	}
	if lookUpErr != nil {
		http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
		fmt.Println("[ERROR] ", lookUpErr)
		return
	}
	a := data.Auth{UserId: userID, Otp: dto.Otp}
	if result, err := a.CheckOTP(); err != nil {
		http.Error(w, "Error Ocurred while verifying OTP", http.StatusInternalServerError)
		fmt.Println("[ERROR] ", err)
		return
	} else if !result {
		http.Error(w, "OTP Does not Match", http.StatusBadRequest)
		return
	}

	if result, err := utils.GenerateJWT(userID); err != nil {
		http.Error(w, "Error Occured during access token generation", http.StatusBadRequest)
		fmt.Println("[ERROR] ", err)
		return
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, data.AccessTokenDTO{AccessToken: result})
	}

}

func LoginViaPhone(w http.ResponseWriter, r *http.Request) {
	var dto data.LoginWithPhoneDTO
	err := utils.DecodeJSONRequest(r, dto)
	if err != nil {
		utils.EncodeJSONResponse(w, http.StatusBadRequest, struct {
			utils.GenericJsonResponseDTO
			Err string `json:"err"`
		}{
			GenericJsonResponseDTO: utils.GenericJsonResponseDTO{
				Message: "Invalid Request Body",
			},
			Err: err.Error(),
		})
		return
	}
	result, userID, lookUpErr := data.UserLookupViaPhone(dto.Phone)

	if !result {
		http.Error(w, "User Does not Exist", http.StatusNotFound)
		return
	}
	if lookUpErr != nil {
		http.Error(w, "Error Ocurred while Lookup", http.StatusInternalServerError)
		fmt.Println("[ERROR] ", lookUpErr)
		return
	}
	a := data.Auth{UserId: userID, Otp: dto.Otp}
	if result, err := a.CheckOTP(); err != nil {
		http.Error(w, "Error Ocurred while verifying OTP", http.StatusInternalServerError)
		fmt.Println("[ERROR] ", err)
		return
	} else if !result {
		http.Error(w, "OTP Does not Match", http.StatusBadRequest)
		return
	}

	if result, err := utils.GenerateJWT(userID); err != nil {
		http.Error(w, "Error Occured during access token generation", http.StatusBadRequest)
		fmt.Println("[ERROR] ", err)
		return
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, data.AccessTokenDTO{AccessToken: result})
	}
}
