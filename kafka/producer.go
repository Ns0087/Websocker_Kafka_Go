package kafka

import (
	"WebsocketExample/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Kafka writer (producer)
var writer *kafka.Writer

func InitProducer(broker, topic string) {
	writer = &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne, // Ensures at least one broker gets the message
		BatchTimeout:           10 * time.Millisecond,
		AllowAutoTopicCreation: true,
	}
}

func PublishMessage(topic string, message *models.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error while serializing message model to json bytes", err)
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(message.Receiver),
		Value: messageBytes,
	})
	if err != nil {
		fmt.Println("Error while writing message to kafka", err)
		return err
	}
	return nil
}

func CloseProducer() {
	err := writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer", err)
	} else {
		fmt.Println("Producer is closed")
	}
}
