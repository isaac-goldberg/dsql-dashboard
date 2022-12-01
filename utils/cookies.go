package utils

import (
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, key string, value string) {
	var cookie = http.Cookie{
		Name:    key,
		Value:   value,
		Expires: time.Now().Add(time.Second * time.Duration(604800)),
	}

	http.SetCookie(w, &cookie)
}

func SetTokensAndIdCookies(w http.ResponseWriter, id string, accessToken string, refreshToken string) {
	SetCookie(w, "t", Encrypt(accessToken))
	SetCookie(w, "r", Encrypt(refreshToken))
	SetCookie(w, "id", id)
}

func DeleteTokenAndIdCookies(w http.ResponseWriter) {
	DeleteCookie(w, "t")
	DeleteCookie(w, "r")
	DeleteCookie(w, "id")
}
func DeleteCookie(w http.ResponseWriter, key string) {
	var cookie = http.Cookie{
		Name:   key,
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
}
