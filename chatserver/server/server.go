package server

import (
	"fmt"

	channel "github.com/socketChat/chatserver/channel"
)

var (
	ErrChannelExist    = fmt.Errorf("Channel exists")
	ErrChannelNotExist = fmt.Errorf("Channel does not exist")
)

type ChanMgr struct {
	addr     string
	channels map[string]*channel.Channel
}

func NewChanMgr(addr string) *ChanMgr {
	return &ChanMgr{
		addr:     addr,
		channels: map[string]*channel.Channel{},
	}
}

func (s *ChanMgr) AddChannel(ch *channel.Channel) error {
	if _, exist := s.channels[ch.ChannelID]; exist {
		return ErrChannelExist
	}
	s.channels[ch.ChannelID] = ch
	return nil
}

func (s *ChanMgr) GetChannel(channelID string) (*channel.Channel, error) {
	ch, exist := s.channels[channelID]
	if !exist {
		return nil, ErrChannelNotExist
	}
	return ch, nil
}
