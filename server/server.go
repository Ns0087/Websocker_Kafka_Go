package server

import (
	"WebsocketExample/kafka"
	"WebsocketExample/websocket"
	"flag"
	"fmt"
	"net/http"
)

// StartServer initializes the WebSocket server
func StartServer() {
	port := flag.String("port", "8081", "server port")
	flag.Parse()

	// Kafka Broker and Topic
	kafkaBroker := "localhost:9092"
	chatTopic := "chat-messages"
	notifyTopic := "notifications"

	// Initialize Kafka Producer
	kafka.InitProducer(kafkaBroker, chatTopic)

	// Set WebSocket to use Kafka publisher (breaking the cycle)
	websocket.PublishToKafka = kafka.PublishMessage

	// Start Kafka Consumers
	go kafka.StartConsumer(kafkaBroker, chatTopic, websocket.SendPrivateMessage)
	go kafka.StartConsumer(kafkaBroker, notifyTopic, websocket.BroadcastNotification)

	// Websocket Handlers
	http.HandleFunc("/ws/chat", websocket.HandleChats)
	http.HandleFunc("/ws/notify", websocket.HandleNotifications)

	serverAddress := fmt.Sprintf(":%s", *port)
	fmt.Println("Server started at " + serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		fmt.Println("Server Error: ", err)
	}
}
