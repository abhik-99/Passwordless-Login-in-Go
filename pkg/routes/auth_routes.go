package routes

import (
	"github.com/abhik-99/passwordless-login/pkg/controllers"

	"github.com/gorilla/mux"
)

// Un-Protected Routes
func RegisterAuthRoutes(r *mux.Router) {

	r.HandleFunc("/signup", controllers.Signup).Methods("POST")

	//This initiates the Email-based authentication,
	r.HandleFunc("/email/{emailId}", controllers.OTPViaEmail).Methods("GET")
	//the Email & OTP can be passed via REQ Body
	r.HandleFunc("/email", controllers.LoginViaEmail).Methods("POST")

	//This initiates the Phone-based authentication,
	r.HandleFunc("/phone/{phoneNo}", controllers.OTPViaPhone).Methods("GET")
	// the Phone number & OTP can be passed via REQ Body
	r.HandleFunc("/phone", controllers.LoginViaPhone).Methods("POST")

}
