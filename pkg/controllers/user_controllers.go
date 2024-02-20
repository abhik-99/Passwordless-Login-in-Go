package controllers

import (
	"log"
	"net/http"

	"github.com/abhik-99/passwordless-login/pkg/data"
	"github.com/abhik-99/passwordless-login/pkg/utils"
	"github.com/gorilla/mux"
)

func GetPublicUsers(w http.ResponseWriter, r *http.Request) {
	if userProfiles, err := data.GetAllPublicUserProfiles(); err != nil {
		http.Error(w, "Error Occurred while fetching users", http.StatusInternalServerError)
		log.Println("ERROR", err)
	} else {
		utils.EncodeJSONResponse(w, http.StatusOK, userProfiles)
	}
}

func GetPublicUserProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if !utils.IsValidObjectID(id) {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}
	user, err := data.GetUserProfileById(id)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		log.Println("ERROR while user profile query", err)
		return
	}
	utils.EncodeJSONResponse(w, http.StatusOK, user)

}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user")
	user, err := data.GetUserProfileById(userID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		log.Println("ERROR while user profile query", err)
		return
	}
	utils.EncodeJSONResponse(w, http.StatusOK, user)

}
