package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jdpolicano/govault/internal/server"
	e "github.com/jdpolicano/govault/internal/server/errors"
)

// Middleware defines a function that wraps an http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies several middleware around a final handler.
func Chain(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// ValidateToken validates the Authorization header and stores the session on the request context.
func ValidateToken(ctx *server.Context) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sess, err := ctx.ValidateTokenHeader(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), server.SessionKey{}, sess))
			next(w, r)
		}
	}
}

// ParseJSONBody parses the request body JSON into a value of type T and stores it on the context.
func ParseJSONBody[T any]() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var body T
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&body); err != nil {
				http.Error(w, e.InvalidRequestBody.Error(), http.StatusBadRequest)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), server.BodyKey{}, body))
			next(w, r)
		}
	}
}

// Logging logs the request method, path and duration using the provided logger.
func Logging(l *log.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			l.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		}
	}
}
