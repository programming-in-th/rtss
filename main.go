package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/programming-in-th/rtss/ws"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id   uint64
	hub  *Hub
	send chan Payload
}

type Payload struct {
	Id     uint64  `json:"id"`
	Groups []Group `json:"groups"`
	Score  int64   `json:"score"`
	Status string  `json:"status"`
}

func SSE(hub *Hub, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	rawId := r.URL.Query().Get("id")
	id, _ := strconv.ParseUint(rawId, 10, 32)

	client := &Client{id: id, hub: hub, send: make(chan Payload, 256)}
	client.hub.register <- client
	defer func() {
		client.hub.unregister <- client
	}()

	timeout := time.After(600 * time.Second)

	for {
		select {
		case message := <-client.send:
			fmt.Fprintf(w, "data: %s\n\n", PayloadToJSONString(message))
			w.(http.Flusher).Flush()

		case <-r.Context().Done():
			return
		case <-timeout:
			return
		}
	}

}

func main() {
	hub := newHub()
	go hub.run()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	wsHost := os.Getenv("WS_HOST")
	if wsHost == "" {
		wsHost = "157.230.244.51:4000"
	}

	wsPath := os.Getenv("WS_PATH")
	if wsPath == "" {
		wsPath = "/socket/websocket"
	}

	go func() {
		u := url.URL{Scheme: "ws", Host: wsHost, Path: wsPath}
		s := &ws.Socket{UrlString: u}

		s.Connect()

		defer s.Connection.Close()

		channel := s.SetChannel("realtime:public:Submission")

		channel.Join()
		channel.On("*", func(data interface{}) {
			d := data.(map[string]interface{})["record"]

			if d != nil {
				rawId := d.(map[string]interface{})["id"].(string)
				id, _ := strconv.ParseUint(rawId, 10, 32)

				raw := d.(map[string]interface{})["groups"]

				if raw.(string) == "unchanged_toast" {
					o := data.(map[string]interface{})["old_record"]
					raw = o.(map[string]interface{})["groups"]
				}

				var groups []Group
				json.Unmarshal([]byte(raw.(string)), &groups)
				score, _ := strconv.ParseInt(d.(map[string]interface{})["score"].(string), 10, 32)

				payload := Payload{Id: id, Groups: groups, Score: score, Status: d.(map[string]interface{})["status"].(string)}

				hub.broadcast <- payload
			}

		})

		s.Listen()
	}()

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		SSE(hub, w, r)
	})
	log.Fatal("HTTP server error: ", http.ListenAndServe(":"+port, nil))

}
