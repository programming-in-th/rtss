package channel

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type PhxJson struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload string      `json:"payload"`
	Ref     interface{} `json:"ref"`
}

type CallBackListener struct {
	Callback func(map[string][]byte)
	Event    string
}

type Channel struct {
	Socket    *websocket.Conn
	Topic     string
	Listeners []CallBackListener
}

func GenJson(topic string) string {
	pl := &PhxJson{
		Topic:   topic,
		Event:   "phx_join",
		Payload: "",
		Ref:     nil,
	}

	data, _ := json.Marshal(pl)

	return string(data)
}

func (c Channel) Join() {

}
