package controllers

import (
	"log"
	"net/http"

	"github.com/abhik-99/passwordless-login/pkg/data"
	"github.com/abhik-99/passwordless-login/pkg/utils"
)

func GetPublicUsers(w http.ResponseWriter, r *http.Request) {
	if userProfiles, err := data.GetAllPublicUserProfiles(); err != nil {
		http.Error(w, "Error Occurred while fetching users", http.StatusInternalServerError)
		log.Println("ERROR", err)
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, userProfiles)
	}
}

func GetPublicUserProfile(w http.ResponseWriter, r *http.Request) {}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {}
