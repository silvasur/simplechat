package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func AcceptWebSock(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	room := vars["chatroom"]
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		var nick string
		if websocket.Message.Receive(ws, &nick) != nil {
			return
		}

		if err := Join(room, nick); err != nil {
			// TODO: report error to client
			return
		}

		usermsgs := make(chan string)
		go func() {
			var s string
			for {
				if websocket.Message.Receive(ws, &s) != nil {
					return
				}

				if s != "" {
					usermsgs <- s
				}
			}
		}()
	}).ServeHTTP(rw, req)
}
