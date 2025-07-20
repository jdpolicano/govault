package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
func ValidateToken(refs *server.ServerRefs) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// check that the header exists
			auth := r.Header.Get("Authorization")
			if len(auth) == 0 {
				http.Error(w, "missing authorization", http.StatusUnauthorized)
				return
			}

			// check that the prefix is correct
			toke, found := strings.CutPrefix(auth, "Bearer govault-")
			if !found {
				http.Error(w, "malformed header", http.StatusUnauthorized)
				return
			}

			// check if the token has expired...
			sess, found := refs.Sessions.Get(toke)
			if sess.Expired() {
				refs.Sessions.Delete(toke)
				http.Error(w, "no such session", http.StatusUnauthorized)
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
