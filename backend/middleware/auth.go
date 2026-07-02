package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
)

func AuthMiddleware(app *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client, err := app.Auth(context.Background())
			if err != nil {
				log.Printf("app.Auth() failed: %v", err)
				http.Error(w, "Auth client error", http.StatusInternalServerError)
				return
			}

			// Extract token from "Authorization: Bearer <token>"
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			token, err := client.VerifyIDToken(context.Background(), tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Pass the User ID to the next handler
			ctx := context.WithValue(r.Context(), "user_id", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
