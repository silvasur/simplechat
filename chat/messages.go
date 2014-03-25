package chat

import (
	"encoding/json"
)

// MsgType describes the purpose of a message
type MsgType int

const (
	MsgChat MsgType = iota // Default
	MsgJoin
	MsgLeave
)

func (mt MsgType) String() string {
	switch mt {
	case MsgChat:
		return "chat"
	case MsgJoin:
		return "join"
	case MsgLeave:
		return "leave"
	}

	return "???"
}

func (mt *MsgType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

// Message represents a message that can be sent to a buddy. The Text field has no meaning, if Type != MsgChat.
type Message struct {
	Type MsgType `json:"type"`
	User string  `json:"user"`
	Text string  `json:"text,omitempty"`
}
