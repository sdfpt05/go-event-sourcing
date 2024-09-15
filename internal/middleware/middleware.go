package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()
		r.Header.Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			requestID,
			time.Since(start),
		)
	}
}

func Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}