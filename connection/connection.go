package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/programming-in-th/rtss/channel"
)

type Socket struct {
	Connection *websocket.Conn
	Connected  bool
	UrlString  url.URL
	Channels   []channel.Channel
}

func (s *Socket) Connect() {

	c, _, err := websocket.DefaultDialer.Dial(s.UrlString.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	s.Connected = true
	s.Connection = c
}

func (s *Socket) SetChannel(topic string) *channel.Channel {
	channel := &channel.Channel{Socket: s.Connection, Topic: topic}
	s.Channels = append(s.Channels, *channel)
	return channel
}

func (s *Socket) Listen() {
	defer s.Connection.Close()

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker.C:
				data, _ := json.Marshal(&channel.PhxJson{Topic: "phoenix", Event: "heartbeat", Payload: "{\"msg\": \"ping\"}", Ref: nil})
				err := s.Connection.WriteMessage(websocket.PingMessage, data)

				if err != nil {
					log.Println("Error:", err)
				}
			}
		}
	}()

	for {
		_, message, err := s.Connection.ReadMessage()

		if err != nil {
			log.Println("Error:", err)
		}

		fmt.Println(message)
	}

}
