package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	"github.com/isaac-goldberg/dsql-dashboard/globals"
	"github.com/isaac-goldberg/dsql-dashboard/middleware"
	"github.com/isaac-goldberg/dsql-dashboard/routes"
)

func main() {
	// load env file
	godotenv.Load()

	// handle asset routes
	cssHandler := http.FileServer(http.Dir("./public/assets/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))

	// gorilla mux router
	muxRouter := mux.NewRouter()
	muxRouter.Use(middleware.ValidateUser)
	routes.AddRootRoutesHandler(muxRouter)
	routes.AddAuthRoutesHandler(muxRouter)

	http.Handle("/", muxRouter)
	fmt.Printf("Server online on port %d\n", globals.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", globals.Port), nil))
}
