package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/takuchi17/term-keeper/configs"
)

type contextKey struct {
	name string
}

var (
	userIDKey   = &contextKey{"userID"}
	userNameKey = &contextKey{"userName"}
)

var jwtSecret = []byte(configs.Config.JWTSecret)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Warn("Authorization header is missing")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// parse the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				slog.Warn("Unexpected signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			slog.Warn("Failed to parse token", "err", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user ID and username from claims
			userID, ok := claims["userid"].(string)
			if !ok {
				slog.Warn("Invalid user ID in token claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userName, ok := claims["username"].(string)
			if !ok {
				slog.Warn("Invalid username in token claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Store user ID and username in request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, userNameKey, userName)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		slog.Warn("Invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIDKey).(string)
	return userId, ok
}

func GetUserName(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(userNameKey).(string)
	return name, ok
}
