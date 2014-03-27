package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/mux"
	"github.com/kch42/simplechat/chat"
	"net/http"
	"time"
)

type JoinResponse struct {
	OK      bool     `json:"ok"`
	Error   string   `json:"error,omitempty"`
	Buddies []string `json:"buddies,omitempty"`
}

const (
	floodstopReset = 1 * time.Second
	floodstopMax   = 3
)

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

		buddy, room, err := chat.Join(roomname, nick)
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

		exit := make(chan struct{})

		floodstop := NewFloodstop(floodstopReset, floodstopMax)
		defer floodstop.Stop()

		go func() {
			var s string
			for {
				if websocket.Message.Receive(ws, &s) != nil {
					break
				}

				if !floodstop.Ask() {
					// User sent too many requests in a short amount of time. Kick them!
					break
				}

				if s == "" {
					continue
				}

				buddy.Say(s)
			}

			exit <- struct{}{}
		}()

		for {
			select {
			case <-exit:
				return
			case m := <-buddy.Receive:
				if send(m) != nil {
					return
				}
			}
		}
	}).ServeHTTP(rw, req)
}
