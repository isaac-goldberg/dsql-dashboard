package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/isaac-goldberg/dsql-dashboard/sessions"
	"github.com/isaac-goldberg/dsql-dashboard/typings"
	"github.com/isaac-goldberg/dsql-dashboard/utils"
)

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := r.Cookie("id")
		accessToken, err2 := r.Cookie("t")
		refreshToken, err3 := r.Cookie("r")

		if err == nil && err2 == nil && err3 == nil {
			// attempt to get existing session
			session := sessions.Get(userId.Value)
			// execute if session exists
			if session.AccessToken != "" {
				// if session expire date exists and if it has expired at this time
				if session.ExpiresAt > 0 && time.Now().Unix()-session.ExpiresAt >= 0 {
					// try using refresh token
					newAT, newRT, err4 := utils.UseRefreshToken(session.RefreshToken)
					// if refresh token worked, set new tokens in session and cookies
					if err4 == nil {
						sessions.Set(w, userId.Value, newAT, newRT, session.User)
						utils.SetTokensAndIdCookies(w, userId.Value, newAT, newRT)
					}
				}

				context.Set(r, typings.RouteLocalsKeys_User{}, session.User)
				next.ServeHTTP(w, r)
				return
				// execute if session does not exist but access token is saved in cookies
				// this will only execute if the app restarted but client data was saved
			} else if accessToken.Value != "" {
				decryptedAT, err4 := utils.Decrypt(accessToken.Value)
				decryptedRT, err5 := utils.Decrypt(refreshToken.Value)

				if err4 == nil && err5 == nil {
					userMap, err6 := utils.GetUser(decryptedAT)

					// if access token still works
					if err6 == nil {
						// check if refresh token works
						newAT, newRT, err7 := utils.UseRefreshToken(decryptedRT)

						// if it does, set request token cookies to the new ones
						if err7 == nil {
							decryptedAT = newAT
							decryptedRT = newRT
						}

						id := userMap["id"].(string)
						sessions.SaveUserMap(w, decryptedAT, decryptedRT, userMap)
						session = sessions.Get(id)
						utils.SetTokensAndIdCookies(w, id, decryptedAT, decryptedRT)

						// set user in request context
						context.Set(r, typings.RouteLocalsKeys_User{}, session.User)
						next.ServeHTTP(w, r)
						return
					}
				}

			}
		}

		// if no user exists, delete cookies and continue
		utils.DeleteTokenAndIdCookies(w)
		next.ServeHTTP(w, r)
	})
}
