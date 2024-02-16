package routes

import (
	"github.com/abhik-99/passwordless-login/pkg/controllers"

	"github.com/gorilla/mux"
)

// Un-Protected Routes
func RegisterAuthRoutes(r *mux.Router) {

	//Because this initiates the authentication, the Email or Phone number can be passed via URL
	r.HandleFunc("/email/{emailId}", controllers.LoginViaEmail).Methods("GET")
	r.HandleFunc("/email/{emailId}/verify-otp", controllers.VerifyEmailLoginOTP).Methods("POST")

	//Because this initiates the authentication, the Email or Phone number can be passed via URL
	r.HandleFunc("/phone/{phoneNo}", controllers.LoginViaPhone).Methods("GET")
	r.HandleFunc("/phone/{phoneNo}/verify-otp", controllers.VerifyPhoneLoginOTP).Methods("POST")
}
