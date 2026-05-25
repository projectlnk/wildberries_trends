package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wildberries-trends/internal/metrics"
	"wildberries-trends/internal/stoplist"
	"wildberries-trends/internal/top"
	"wildberries-trends/internal/window"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type TopResponse struct {
	Top []string `json:"top"`
}

func NewRouter(win *window.SlidingWindow, sl *stoplist.StopList) http.Handler {
	mux := http.NewServeMux()

	// GET /top?limit=N
	mux.HandleFunc("GET /top", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			metrics.TopRequestDuration.Observe(time.Since(start).Seconds())
		}()
		metrics.TopRequestsCount.Inc()

		limitStr := r.URL.Query().Get("limit")
		limit := 10
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}
		counts := win.GetAllCounts()
		// Фильтруем стоп-лист
		filtered := make(map[string]int64)
		for q, c := range counts {
			if !sl.Contains(q) {
				filtered[q] = c
			}
		}
		topQueries := top.GetTopN(filtered, limit)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TopResponse{Top: topQueries})
	})

	// GET /stoplist
	mux.HandleFunc("GET /stoplist", func(w http.ResponseWriter, r *http.Request) {
		words := sl.List()
		json.NewEncoder(w).Encode(map[string][]string{"stoplist": words})
	})

	// POST /stoplist
	mux.HandleFunc("POST /stoplist", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Word string `json:"word"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Word == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		sl.Add(strings.ToLower(req.Word))
		json.NewEncoder(w).Encode(map[string]string{"status": "added", "word": req.Word})
	})

	// DELETE /stoplist/{word}
	mux.HandleFunc("DELETE /stoplist/{word}", func(w http.ResponseWriter, r *http.Request) {
		word := r.PathValue("word")
		if word == "" {
			http.Error(w, "Word required", http.StatusBadRequest)
			return
		}
		sl.Remove(strings.ToLower(word))
		json.NewEncoder(w).Encode(map[string]string{"status": "removed", "word": word})
	})

	// Метрики Prometheus
	mux.Handle("GET /metrics", promhttp.Handler())

	return mux
}
