package middleware

import (
	"net/http"
	"strings"

	"github.com/abhik-99/passwordless-login/pkg/utils"
)

func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(tokenString, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString = splitToken[1]

		tokenClaims, err := utils.ValidateJWT(tokenString)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userId, _ := tokenClaims.GetSubject()
		r.Header.Set("user", userId)

		next.ServeHTTP(w, r)
	})
}
