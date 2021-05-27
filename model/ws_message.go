package model

import (
	"encoding/json"
	"log"
)

type ReceivedMessage struct {
	Action string `json:"action"`
	Room   string `json:"room"`
	Message *string `json:"message"`
}

type WebsocketMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (message *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return encoding
}
