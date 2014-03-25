package chat

import (
	"time"
)

type Buddy struct {
	Nick    string
	Receive chan Message
	room    *Room
}

func NewBuddy(nick string, room *Room) *Buddy {
	return &Buddy{
		Nick:    nick,
		Receive: make(chan Message),
		room:    room,
	}
}

func (b *Buddy) Leave() {
	b.room.Leave(b.Nick)
}

func (b *Buddy) Push(msg Message) {
	go func() {
		select {
		case b.Receive <- msg:
		case <-time.Tick(time.Millisecond * 100):
		}
	}()
}

// Say sends a text as a chat message of this user to the connected room.
func (b *Buddy) Say(text string) {
	b.room.Messages <- Message{
		Type: MsgChat,
		User: b.Nick,
		Text: text,
	}
}
