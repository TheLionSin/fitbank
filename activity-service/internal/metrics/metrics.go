package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActivitiesCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fitbank_activities_created_total",
		Help: "The total number of activities created",
	})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fitbank_request_duration_seconds",
		Help:    "Histogram of request durations",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
	}, []string{"method", "path"})
)
