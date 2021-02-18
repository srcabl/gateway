package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// InjectCors if the middleware handler for CORS
func InjectCors() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowCredentials: true,
		//Debug:            true,
	}).Handler
}
