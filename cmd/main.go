package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abhik-99/passwordless-login/pkg/config"
	"github.com/abhik-99/passwordless-login/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {
	config.Connect()
	defer config.Disconnect()
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	routes.RegisterAuthRoutes(authRouter)

	userRouter := router.PathPrefix("/user").Subrouter()
	routes.RegisterUserRoutes(userRouter)

	fmt.Println("Starting on PORT 3000")
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", router))
}
