package ws

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Socket struct {
	Connection *websocket.Conn
	Connected  bool
	UrlString  url.URL
	Channels   []*Channel
}

type Reply struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
	Ref     interface{} `json:"ref"`
	Topic   string      `json:"topic"`
}

func (s *Socket) Connect() {
	c, _, err := websocket.DefaultDialer.Dial(s.UrlString.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Printf("connecting to %s", s.UrlString.String())

	s.Connected = true
	s.Connection = c
}

func (s *Socket) SetChannel(topic string) *Channel {
	channel := &Channel{Socket: s.Connection, Topic: topic}
	s.Channels = append(s.Channels, channel)
	return channel
}

func (s *Socket) Listen() {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := s.Connection.ReadMessage()

			var data Reply
			json.Unmarshal(message, &data)

			for _, v := range s.Channels {

				if v.Topic == data.Topic {
					for _, v2 := range v.Listeners {
						if v2.Event == data.Event || v2.Event == "*" {
							v2.Callback(data.Payload)
						}
					}
				}
			}

			if err != nil {
				log.Println("Error:", err)
			}

		}
	}()

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			data, _ := json.Marshal(&PhxJson{Topic: "phoenix", Event: "heartbeat", Payload: "{\"msg\": \"ping\"}", Ref: nil})
			err := s.Connection.WriteMessage(websocket.TextMessage, data)

			if err != nil {
				log.Println("Error:", err)
			}
		}
	}

}
