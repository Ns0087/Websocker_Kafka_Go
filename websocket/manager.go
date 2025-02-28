package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	Conn   *websocket.Conn
	UserID string
}

type ConnectionManager struct {
	chatClients         map[string]*Connection // Chat users
	notificationClients map[string]*Connection // Notification users
	mu                  sync.Mutex             // Mutex for concurrency safety
}

// Manager Global connection manager instance
var Manager = ConnectionManager{
	chatClients:         make(map[string]*Connection),
	notificationClients: make(map[string]*Connection),
}

// AddClient Add client to connection pool
func (cm *ConnectionManager) AddClient(userID string, clientType string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	switch clientType {
	case "chat":
		cm.chatClients[userID] = &Connection{
			UserID: userID,
			Conn:   conn,
		}
		fmt.Println("Chat Client Added: ", userID)

	case "notification":
		cm.notificationClients[userID] = &Connection{
			UserID: userID,
			Conn:   conn,
		}
		fmt.Println("Notification Client Added: ", userID)

	default:
		panic("Unknown clientType")
	}

}

// RemoveClient Remove client from the connection poll
func (cm *ConnectionManager) RemoveClient(userID string, clientType string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	switch clientType {
	case "chat":
		if conn, exists := cm.chatClients[userID]; exists {
			err := conn.Conn.Close()
			if err != nil {
				fmt.Println(err)
			}
			delete(cm.chatClients, userID)
			fmt.Println("Chat Client removed: ", userID)
		}
	case "notification":
		if conn, exists := cm.notificationClients[userID]; exists {
			err := conn.Conn.Close()
			if err != nil {
				fmt.Println(err)
			}
			delete(cm.notificationClients, userID)
			fmt.Println("Notification Client removed: ", userID)
		}
	default:
		panic("Unknown clientType")
	}
}
