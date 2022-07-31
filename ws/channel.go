package ws

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
	Callback func(interface{})
	Event    string
}

type Channel struct {
	Socket    *websocket.Conn
	Topic     string
	Listeners []CallBackListener
}

func GenJson(topic string) []byte {
	pl := &PhxJson{
		Topic:   topic,
		Event:   "phx_join",
		Payload: "",
		Ref:     nil,
	}

	data, _ := json.Marshal(pl)

	return data
}

func (c *Channel) Join() {
	j := GenJson(c.Topic)
	c.Socket.WriteMessage(websocket.TextMessage, j)
}

func (c *Channel) On(event string, callback func(interface{})) {
	c.Listeners = append(c.Listeners, CallBackListener{Callback: callback, Event: event})
}

func (c *Channel) Off(event string) {
	for i := range c.Listeners {
		if c.Listeners[i].Event == event {
			c.Listeners = append(c.Listeners[:i], c.Listeners[i+1:]...)
		}
	}
}
