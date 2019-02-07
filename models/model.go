package model

type MsgType int

const (
	Join MsgType = iota
	Text
	Leave
)

type Message struct {
	Type      MsgType `json:"type"`
	OpenID    string  `json:"openID"`
	ChannelID string  `json:"channelID"`
	Text      string  `json:"Text"`
}

type ID string
