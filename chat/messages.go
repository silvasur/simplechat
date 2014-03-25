package chat

import (
	"encoding/json"
	"errors"
)

type MsgType int

const (
	MsgChat MsgType = iota // Default
	MsgJoin
	MsgLeave
)

func (mt *MsgType) MarshalJSON() ([]byte, error) {
	switch *mt {
	case MsgChat:
		return json.Marshal("chat")
	case MsgJoin:
		return json.Marshal("join")
	case MsgLeave:
		return json.Marshal("leave")
	}

	return nil, errors.New("Unknown message type")
}

type Message struct {
	Type MsgType `json:"type"`
	User string  `json:"user"`
	Text string  `json:"text,omitempty"`
}
