package main

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func wsEcho(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Println("Failed to accept websocket:", err)
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	ctx := r.Context()
	for {
		var msg any
		if err := wsjson.Read(ctx, c, &msg); err != nil {
			log.Println("read:", err)
			return
		}

		if err := wsjson.Write(ctx, c, msg); err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsEcho)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK. connect: websocat ws://localhost:8080/ws"))
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
