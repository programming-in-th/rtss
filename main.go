package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// data variable
var msg chan Group

func SSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	msg = make(chan Group)

	defer func() {
		close(msg)
		msg = nil
	}()

	timeout := time.After(600 * time.Second)

	for {
		select {
		case group := <-msg:
			fmt.Fprintf(w, "data: %s\n\n", GroupToJSONString(group))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
		case <-timeout:
			return
		}
	}

}

func main() {
	http.HandleFunc("/stream", SSE)

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:8080", nil))
}
