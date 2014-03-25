package main

type Buddy struct {
	Nick    string
	Receive chan Message
	room    *Room
}

func NewBuddy(nick string, room *Room) *Buddy {
	// TODO: Implement me!
}

func (b *Buddy) Leave() {
	b.room.Leave(b.Nick)
}
