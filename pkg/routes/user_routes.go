package routes

import (
	"passwordless-login/pkg/controllers"

	"github.com/gorilla/mux"
)

// Protected Routes
func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/", controllers.GetPublicUsers).Methods("GET")
	r.HandleFunc("/{id}", controllers.GetPublicUserProfile).Methods("GET")
	r.HandleFunc("/profile", controllers.GetUserProfile).Methods("GET")
}
