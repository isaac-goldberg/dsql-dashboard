package sessions

import (
	"net/http"
	"time"

	"github.com/isaac-goldberg/dsql-dashboard/typings"
)

const RefreshTokenThreshold int64 = 30

type Session struct {
	AccessToken  string
	RefreshToken string
	ResWriter    http.ResponseWriter
	User         typings.UserData
	ExpiresAt    int64
}

var sessions = make(map[string]*Session)

func Has(userId string) bool {
	if _, exists := sessions[userId]; exists {
		return true
	}
	return false
}

func Get(userId string) Session {
	s := sessions[userId]
	if s != nil {
		return *s
	}
	return Session{}
}

func Set(w http.ResponseWriter, userId string, accessToken string, refreshToken string, user typings.UserData) {
	sessions[userId] = &Session{
		ResWriter:    w,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    time.Now().Unix() + RefreshTokenThreshold,
	}
}

func SaveUserMap(w http.ResponseWriter, accessToken string, refreshToken string, userMap map[string]interface{}) {
	id := userMap["id"].(string)
	iconHash := userMap["avatar"].(string)
	var user = typings.UserData{
		UserId:        id,
		Username:      userMap["username"].(string),
		Discriminator: userMap["discriminator"].(string),
		IconURL:       "https://cdn.discordapp.com/avatars/" + id + "/" + iconHash + ".png",
	}

	Set(w, id, accessToken, refreshToken, user)
}

func Delete(userId string) {
	delete(sessions, userId)
}
