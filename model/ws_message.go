package model

import (
	"encoding/json"
	"log"
)

// ReceivedMessage represents a received websocket message
type ReceivedMessage struct {
	Action  string  `json:"action"`
	Room    string  `json:"room"`
	Message *string `json:"message"`
}

// WebsocketMessage represents an emitted message
type WebsocketMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// Encode turns the message into a byte array
func (message *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return encoding
}
