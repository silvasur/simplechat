package main

type Buddy struct {
	Nick    string
	Receive chan Message
}

func NewBuddy(nick string, room *Room) *Buddy {
	// TODO: Implement me!
}
