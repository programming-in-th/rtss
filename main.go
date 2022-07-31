package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/programming-in-th/rtss/connection"
	"github.com/programming-in-th/rtss/group"
)

// data variable
var msg chan group.Group

func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	msg = make(chan group.Group)

	defer func() {
		close(msg)
		msg = nil
	}()

	timeout := time.After(600 * time.Second)

	for {
		select {
		case g := <-msg:
			fmt.Fprintf(w, "data: %s\n\n", group.GroupToJSONString(g))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
		case <-timeout:
			return
		}
	}

}

func main() {
	// http.HandleFunc("/stream", SSE)

	// log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:8080", nil))

	// @TODO pass message to SSE
	u := url.URL{Scheme: "ws", Host: "157.230.244.51:4000", Path: "/socket/websocket"}

	s := &connection.Socket{UrlString: u}

	s.Connect()

	channel := s.SetChannel("realtime:public:Submission")

	channel.Join()
	channel.On("*", func(data interface{}) {
		fmt.Println(data)
	})

	s.Listen()

}
