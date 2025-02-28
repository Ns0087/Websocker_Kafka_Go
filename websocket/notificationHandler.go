package websocket

import (
	"WebsocketExample/models"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// HandleNotifications manages WebSocket connections for notifications
func HandleNotifications(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Websoket upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Extract user ID from query param (example)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		_ = conn.Close()
		return
	}

	// Add user to notification clients
	Manager.AddClient(userID, "notification", conn)
	fmt.Println("New Notification Connection Established")

	// Keep connection open
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			fmt.Println(err)
			Manager.RemoveClient(userID, "notification")
			break
		}
	}
}

// BroadcastNotification Broadcast a notification to only Notification clients
func BroadcastNotification(msg models.Message, _ bool) {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	// Send notification to all connected notification clients
	for _, conn := range Manager.notificationClients {
		err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Message))
		if err != nil {
			fmt.Println("Error Sending Notification: ", err)
			Manager.RemoveClient(conn.UserID, "notification")
		}
	}

	fmt.Println("Broadcast message", msg)
}

// Simulates periodic notifications
//func sendNotifications() {
//	for {
//		time.Sleep(5 * time.Second) // Simulate periodic notifications
//		message := fmt.Sprintf("New notification at %s", time.Now().Format("15:04:05"))
//		BroadcastNotification(models.Message{Message: message}, true)
//	}
//}
