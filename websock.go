package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
	"net/http"
)

type JoinResponse struct {
	OK      bool     `json:"ok"`
	Error   string   `json:"error,omitempty"`
	Buddies []string `json:"buddies,omitempty"`
}

func AcceptWebSock(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	roomname := vars["chatroom"]
	websocket.Handler(func(ws *websocket.Conn) {
		send := func(v interface{}) error { return websocket.JSON.Send(ws, v) }

		defer ws.Close()

		var nick string
		if websocket.Message.Receive(ws, &nick) != nil {
			return
		}

		buddy, room, err := Join(roomname, nick)
		if err != nil {
			send(JoinResponse{
				OK:    false,
				Error: err.Error(),
			})
			return
		}
		defer buddy.Leave()

		if send(JoinResponse{
			OK:      true,
			Buddies: room.ListBuddies(),
		}) != nil {
			return
		}

		go func() {
			var s string
			for {
				if websocket.Message.Receive(ws, &s) != nil {
					return
				}

				if s == "" {
					continue
				}

				buddy.Say(s)
			}
		}()

		for m := range buddy.Receive {
			if send(m) != nil {
				return
			}
		}
	}).ServeHTTP(rw, req)
}
