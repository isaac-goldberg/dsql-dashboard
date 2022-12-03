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
	"github.com/isaac-goldberg/dsql-dashboard/utils"
)

func main() {
	// load env file
	godotenv.Load()

	// handle asset routes
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/assets/js/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./public/assets/fonts/"))))
	http.Handle("/css/", utils.NoCache(http.StripPrefix("/css/", http.FileServer(http.Dir("./public/assets/css/")))))

	// gorilla mux router
	muxRouter := mux.NewRouter()
	muxRouter.Use(middleware.ValidateUser)
	routes.AddRootRoutesHandler(muxRouter)
	routes.AddAuthRoutesHandler(muxRouter)

	http.Handle("/", muxRouter)
	fmt.Printf("Server online on port %d\n", globals.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", globals.Port), nil))
}
