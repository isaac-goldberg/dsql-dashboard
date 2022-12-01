package routes

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/isaac-goldberg/dsql-dashboard/typings"
	"github.com/isaac-goldberg/dsql-dashboard/utils"
)

func AddRootRoutesHandler(router *mux.Router) {
	router.HandleFunc("/", homeRouter)
}

func homeRouter(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, typings.RouteLocalsKeys_User{})

	templateData := typings.TemplatingData{Title: "Home"}
	if user != nil {
		templateData.User = user.(typings.UserData)
	}

	utils.ParsedTemplates.ExecuteTemplate(w, "home.html", templateData)
}
