package chat

import (
	"time"
)

// Buddy represents a user that participates in a chat.
// Read from the Receive channel to get incoming messages
type Buddy struct {
	Nick    string
	Receive chan Message
	room    *Room
}

func newBuddy(nick string, room *Room) *Buddy {
	return &Buddy{
		Nick:    nick,
		Receive: make(chan Message),
		room:    room,
	}
}

// Leave will remove the buddy from the room.
func (b *Buddy) Leave() {
	b.room.leave(b.Nick)
}

// Push pushes a message to the buddies Receive channel
func (b *Buddy) Push(msg Message) {
	go func() {
		t := time.NewTicker(time.Millisecond * 100)
		defer t.Stop()
		select {
		case b.Receive <- msg:
		case <-t.C:
		}
	}()
}

// Say sends a text as a chat message of this user to the connected room.
func (b *Buddy) Say(text string) {
	b.room.messages <- Message{
		Type: MsgChat,
		User: b.Nick,
		Text: text,
	}
}
