package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"wildberries-trends/internal/api"
	"wildberries-trends/internal/consumer"
	"wildberries-trends/internal/stoplist"
	"wildberries-trends/internal/window"
	"wildberries-trends/pkg/models"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "search-events"

	win := window.NewSlidingWindow(5*time.Minute, 5*time.Second)
	sl := stoplist.New()

	handler := func(ev *models.SearchEvent) {
		win.Add(ev.Query)
	}

	consumerInst, err := consumer.NewConsumer(brokers, topic, handler)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	router := api.NewRouter(win, sl)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		log.Println("Starting API server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP error: %v", err)
		}
	}()

	log.Println("Starting Kafka consumer...")
	if err := consumerInst.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
}
