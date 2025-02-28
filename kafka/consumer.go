package kafka

import (
	"WebsocketExample/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// StartConsumer Consume messages from Kafka and broadcast them via WebSocket
func StartConsumer(broker, topic string, handler MessageHandler) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		StartOffset:    kafka.LastOffset, // Start reading only new messages
		CommitInterval: 0,                // Commit messages immediately
		MaxBytes:       10e6,
	})
	defer reader.Close()

	log.Println("Kafka Consumer started for topic:", topic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading Kafka message:", err)
			continue
		}

		// Parse message (format: "receiverID:message")
		message := models.Message{}
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Println("Error unmarshalling Kafka message:", err)
			continue
		}

		fmt.Printf("Message from %s to %s: %s\n", message.Sender, message.Receiver, message.Message)

		// Send message to WebSocket clients
		handler(message, true)
	}
}
