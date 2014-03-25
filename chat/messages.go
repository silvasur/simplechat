package chat

import (
	"encoding/json"
	"errors"
)

// MsgType describes the purpose of a message
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

// Message represents a message that can be sent to a buddy. The Text field has no meaning, if Type != MsgChat.
type Message struct {
	Type MsgType `json:"type"`
	User string  `json:"user"`
	Text string  `json:"text,omitempty"`
}
