package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/isaac-goldberg/dsql-dashboard/sessions"
	"github.com/isaac-goldberg/dsql-dashboard/typings"
	"github.com/isaac-goldberg/dsql-dashboard/utils"
)

func AddAuthRoutesHandler(router *mux.Router) {
	router.HandleFunc("/login", loginRouter)
	router.HandleFunc("/logout", logoutRouter)
	router.HandleFunc("/auth", authRouter)
}

func loginRouter(w http.ResponseWriter, r *http.Request) {
	loginUri := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=919344970551414844&redirect_uri=%s/auth&response_type=code&scope=guilds identify&prompt=none", url.QueryEscape(os.Getenv("HOSTNAME")))

	http.Redirect(w, r, loginUri, http.StatusTemporaryRedirect)
}

func authRouter(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	accessToken, refreshToken, err := utils.GetAccessTokens(code)
	userMap, err2 := utils.GetUser(accessToken)

	if err == nil && err2 == nil {
		sessions.SaveUserMap(w, accessToken, refreshToken, userMap)
		utils.SetTokensAndIdCookies(w, userMap["id"].(string), accessToken, refreshToken)
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func logoutRouter(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, typings.UserData{})

	if user != nil {
		u := user.(typings.UserData)
		session := sessions.Get(u.UserId)
		if session.User.UserId != "" {
			sessions.Delete(u.UserId)
		}
	}

	utils.DeleteTokenAndIdCookies(w)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
