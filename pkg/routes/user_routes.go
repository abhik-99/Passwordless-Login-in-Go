package routes

import (
	"github.com/abhik-99/passwordless-login/pkg/controllers"

	"github.com/gorilla/mux"
)

// Protected Routes
func RegisterUserRoutes(r *mux.Router) {
	// r.Use(middleware.ValidateTokenMiddleware)
	r.HandleFunc("/", controllers.GetPublicUsers).Methods("GET")
	r.HandleFunc("/{id}", controllers.GetPublicUserProfile).Methods("GET")
	r.HandleFunc("/profile", controllers.GetUserProfile).Methods("GET")
}
