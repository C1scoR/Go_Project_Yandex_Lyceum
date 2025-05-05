package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	emailKey  contextKey = "email"
	roleKey   contextKey = "role"
)

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "В заголовке не было передано токена для авторизации", http.StatusBadRequest)
				log.Println("middleware/AuthMiddleware(): В заголовке не передали токена для авторизации")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Инвалидные данные для авторизации", http.StatusUnauthorized)
				log.Println("middleware/AuthMiddleware(): Инвалидные данные для авторизации")
				return
			}
			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Инвалидный токен", http.StatusUnauthorized)
				log.Println("middleware/AuthMiddleware(): Инвалидный токен для авторизации: ", err)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Неверные поля токена", http.StatusUnauthorized)
				log.Println("middleware/AuthMiddleware(): Неверные поля токена")
				return
			}
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					http.Error(w, "Токен устарел", http.StatusUnauthorized)
					log.Println("middleware/AuthMiddleware(): Токен устарел")
					return
				}
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims["user_id"])
			ctx = context.WithValue(ctx, emailKey, claims["email"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RateLimiter middleware для предотвращения атак методом подбора
func RateLimiterMiddleware() func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
