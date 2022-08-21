package main

import (
	"fmt"
	"log"
	"net/http"
	"passwordless-login/pkg/config"
	"passwordless-login/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.Connect()
	defer config.Disconnect()
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	routes.RegisterLoginRoutes(authRouter)

	logoutRouter := router.PathPrefix("/de-auth").Subrouter()
	routes.RegisterLogoutRoute(logoutRouter)

	userRouter := router.PathPrefix("/user").Subrouter()
	routes.RegisterUserRoutes(userRouter)

	fmt.Println("Starting on PORT 3000")
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", router))
}
