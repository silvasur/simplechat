package chat

import (
	"errors"
)

// Room represents a chatroom.
type Room struct {
	messages chan Message
	buddies  map[string]*Buddy
}

func newRoom() (r *Room) {
	r = new(Room)
	r.messages = make(chan Message)
	r.buddies = make(map[string]*Buddy)
	go r.broadcast()
	return
}

func (r *Room) leave(nick string) {
	if _, ok := r.buddies[nick]; !ok {
		return
	}

	delete(r.buddies, nick)
	if len(r.buddies) == 0 {
		close(r.messages)
	} else {
		r.messages <- Message{
			Type: MsgLeave,
			User: nick,
		}
	}
}

func (r *Room) broadcast() {
	for m := range r.messages {
		for _, buddy := range r.buddies {
			buddy.Push(m)
		}
	}
}

// ListBuddies returns a list of nicknames of the connected buddies.
func (r *Room) ListBuddies() (buddies []string) {
	for nick := range r.buddies {
		buddies = append(buddies, nick)
	}
	return
}

var (
	rooms   map[string]*Room
	perroom int
)

// InitRooms initializes the internal rooms variable. Use this, before calling Join.
func InitRooms(room_limit int) {
	rooms = make(map[string]*Room)
	perroom = room_limit
}

// Join joins a buddy to a room. The room will be created, if it doesn't exist.
func Join(room, nick string) (*Buddy, *Room, error) {
	r, ok := rooms[room]
	if !ok {
		r = newRoom()
		rooms[room] = r
	}

	if _, there := r.buddies[nick]; there {
		return nil, nil, errors.New("Nickname is already in use")
	}

	if len(r.buddies) >= perroom {
		return nil, nil, errors.New("Room is full")
	}

	r.messages <- Message{
		Type: MsgJoin,
		User: nick,
	}

	b := newBuddy(nick, r)
	r.buddies[nick] = b
	return b, r, nil
}
