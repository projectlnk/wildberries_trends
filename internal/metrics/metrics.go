package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TotalEvents счётчик обработанных поисковых событий
	TotalEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_events_total",
		Help: "Total number of processed search events",
	})
	// TopRequestsCount счётчик запросов к /top
	TopRequestsCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "top_requests_total",
		Help: "Total number of /top requests",
	})
	// TopRequestDuration гистограмма времени ответа /top
	TopRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "top_request_duration_seconds",
		Help:    "Duration of /top requests",
		Buckets: prometheus.DefBuckets,
	})
	// UniqueQueriesInWindow количество уникальных запросов в окне
	UniqueQueriesInWindow = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "unique_queries_in_window",
		Help: "Number of unique queries in the current sliding window",
	})
)
