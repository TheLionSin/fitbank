package middleware

import (
	"fitbank/activity-service/internal/metrics"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Received request %s %s", r.Method, r.URL.Path)
		start := time.Now()

		// Передаем управление следующему хендлеру в цепочке
		next.ServeHTTP(w, r)

		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		log.Printf("METHOD: %s | PATH: %s | DURATION: %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
