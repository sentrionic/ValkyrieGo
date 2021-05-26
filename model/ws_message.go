package model

import (
	"encoding/json"
	"log"
)

type ReceivedMessage struct {
	Action string `json:"action"`
	Room   string `json:"room"`
}

type WebsocketMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func (message *WebsocketMessage) encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return encoding
}
