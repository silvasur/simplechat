package main

import (
	"errors"
)

type Room struct {
	Messages chan Message
	Buddies  map[string]Buddy
}

func NewRoom() (r *Room) {
	r = new(Room)
	r.Messages = make(chan Message)
	r.Buddies = make(map[string]Buddy)
	go r.Broadcast()
}

func (r *Room) Leave(nick string) {
	if _, ok := r.Buddies[nick]; !ok {
		return
	}

	delete(r.Buddies[nick])
	if len(r.Buddies) == 0 {
		close(r.Messages)
	} else {
		r.Messages <- Message{
			Type: MsgLeave,
			From: nick,
		}
	}
}

func (r *Room) Broadcast() {
	for m := range r.Messages {
		for _, buddy := range r.Buddies {
			buddy.Receive <- m // TODO: What happens when this locks?
		}
	}
}

var rooms = make(map[string]Room)

func Join(room, nick string) (*Buddy, *Room, error) {
	r, ok := rooms[room]
	if !ok {
		r = NewRoom()
		rooms[room] = r
	}

	if _, there := r.Buddies[nick]; there {
		return nil, room, errors.New("Nickname is already in use")
	}

	if len(r.Buddies) >= *perroom {
		return nil, room, errors.New("Room is full")
	}

	r.Messages <- Message{
		Type: MsgJoin,
		From: nick,
	}

	b := NewBuddy(nick, r)
	r.Buddies[nick] = b
	return b, nil
}
