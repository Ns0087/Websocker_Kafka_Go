package kafka

import "WebsocketExample/models"

// MessageHandler defines a function signature to handle messages dynamically
type MessageHandler func(message models.Message, fromKafka bool)
