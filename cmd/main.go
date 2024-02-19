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
	defer config.Disconnect()
	router := mux.NewRouter()
	router.StrictSlash(true)

	authRouter := router.PathPrefix("/auth").Subrouter()
	routes.RegisterAuthRoutes(authRouter)

	userRouter := router.PathPrefix("/user").Subrouter()
	routes.RegisterUserRoutes(userRouter)

	fmt.Println("Started on PORT 3000")
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":3000", router))
}
