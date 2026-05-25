package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

// SearchEvent структура, совпадающая с pkg/models/event.go
type SearchEvent struct {
	Query     string `json:"query"`
	Timestamp int64  `json:"timestamp"`
	UserID    string `json:"user_id"`
	IP        string `json:"ip"`
}

var possibleQueries = []string{
	"iphone", "case", "headphones", "watch", "tab",
	"laptop", "mouse", "keyboard", "monitor", "speakers",
}

func main() {
	topic := "search-events"
	brokers := []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Metadata.AllowAutoTopicCreation = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(200 * time.Millisecond) // 5 сообщений в секунду
	defer ticker.Stop()

	log.Println("Producer started. Sending messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Producer stopped.")
			return
		case <-ticker.C:
			query := possibleQueries[rng.Intn(len(possibleQueries))]
			// Имитируем накрутку: каждое 5-е сообщение – "iphone"
			if rng.Intn(5) == 0 {
				query = "iphone"
			}

			event := SearchEvent{
				Query:     query,
				Timestamp: time.Now().Unix(),
				UserID:    fmt.Sprintf("user_%d", rng.Intn(100)),
				IP:        fmt.Sprintf("192.168.%d.%d", rng.Intn(256), rng.Intn(256)),
			}

			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("Marshal error: %v", err)
				continue
			}

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(data),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Send error: %v", err)
			} else {
				log.Printf("Sent: query=%q, partition=%d, offset=%d", query, partition, offset)
			}
		}
	}
}
