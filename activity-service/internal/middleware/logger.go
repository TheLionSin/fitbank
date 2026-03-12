package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Передаем управление следующему хендлеру в цепочке
		next.ServeHTTP(w, r)

		log.Printf("METHOD: %s | PATH: %s | DURATION: %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
