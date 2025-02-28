package websocket

import (
	"WebsocketExample/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// External function to publish messages (set dynamically to avoid cyclic import)
var PublishToKafka func(topic string, message *models.Message) error

// HandleChats manages WebSocket connections for chats
func HandleChats(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query param (example)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Websoket upgrade failed: ", err)
		return
	}

	Manager.AddClient(userID, "chat", conn)

	// Listen for incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			Manager.RemoveClient(userID, "chat")
			break
		}
		fmt.Printf("Message from %s: %s\n", userID, string(msg))

		// Convert JSON message to struct
		message := &models.Message{}
		err = json.Unmarshal(msg, message)
		if err != nil {
			fmt.Println("Json Unmarshal error: ", err)
		}

		// Publish message to Kafka
		if PublishToKafka != nil {
			err = PublishToKafka("chat-messages", message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Send private messages
func SendPrivateMessage(message models.Message, fromKafka bool) {
	if conn, exists := Manager.chatClients[message.Receiver]; exists {
		// Send message directly if user is connected
		err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(message.Message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			conn.Conn.Close()
			Manager.RemoveClient(conn.UserID, "chat")
		}
	}
	if conn, exists := Manager.chatClients[message.Sender]; exists {
		// Send message directly if user is connected
		err := conn.Conn.WriteMessage(websocket.TextMessage, []byte("[ACK] "+message.Message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			conn.Conn.Close()
			Manager.RemoveClient(conn.UserID, "chat")
		}
	} else {
		// no need to send again
		if fromKafka {
			fmt.Println("User", message.Receiver, "is still not connected. Dropping message to avoid loop.")
			return
		}
		// The receiver is not connected to this server instance, but Kafka will route the message
		fmt.Println("User", message.Receiver, "is not connected to this server instance. ending via Kafka.")
		if PublishToKafka != nil {
			err := PublishToKafka("chat-messages", &message)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// AcknowledgeMessage sends an explicit acknowledgment from receiver to sender
func AcknowledgeMessage(message models.Message) {
	ackMessage := models.Message{
		Sender:   message.Receiver, // ACK comes from receiver
		Receiver: message.Sender,   // Sent back to original sender
		Message:  "[DELIVERED] " + message.Message,
	}
	err := PublishToKafka("chat-acknowledgments", &ackMessage)
	if err != nil {
		fmt.Println(err)
	}
}
