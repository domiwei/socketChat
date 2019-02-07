package server

import (
	"fmt"

	channel "github.com/socketChat/chatserver/channel"
)

var (
	ErrChannelExist = fmt.Errorf("Channel exists")
)

type ChanMgr struct {
	addr     string
	Channels map[string]*channel.Channel
}

func NewChanMgr(addr string) *ChanMgr {
	return &ChanMgr{
		addr:     addr,
		Channels: map[string]*channel.Channel{},
	}
}

func (s *ChanMgr) AddChannel(ch *channel.Channel) error {
	if _, exist := s.Channels[ch.ChannelID]; exist {
		return ErrChannelExist
	}
	s.Channels[ch.ChannelID] = ch
	return nil
}
