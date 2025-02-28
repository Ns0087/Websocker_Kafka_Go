package models

type Message struct {
	Sender   string `json:"sender,omitempty"`
	Receiver string `json:"receiver,omitempty"`
	Message  string `json:"message,omitempty"`
}
