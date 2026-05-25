package consumer

import (
	"context"
	"encoding/json"
	"log"
	"wildberries-trends/internal/metrics"
	"wildberries-trends/pkg/models"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	handler  func(*models.SearchEvent)
}

// NewConsumer создаёт нового потребителя Kafka
func NewConsumer(brokers []string, topic string, handler func(*models.SearchEvent)) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // читаем с самого начала (если нужно)

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}, nil
}

// Start запускает прослушивание всех партиций топика
func (c *Consumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					var event models.SearchEvent
					if err := json.Unmarshal(msg.Value, &event); err != nil {
						log.Printf("Failed to unmarshal message: %v", err)
						continue
					}

					metrics.TotalEvents.Inc()
					c.handler(&event)
				case err := <-pc.Errors():
					log.Printf("Consumer error on partition %d: %v", partition, err)
				case <-ctx.Done():
					pc.AsyncClose()
					return
				}
			}
		}(pc)
	}
	<-ctx.Done()
	return nil
}
