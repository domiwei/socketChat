package server

import (
	"fmt"

	channel "github.com/socketChat/chatserver/channel"
)

var (
	ErrChannelExist = fmt.Errorf("Channel exists")
)

type Server struct {
	addr     string
	Channels map[string]*channel.Channel
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		Channels: map[string]*channel.Channel{},
	}
}

func (s *Server) AddChannel(ch *channel.Channel) error {
	if _, exist := s.Channels[ch.ChannelID]; exist {
		return ErrChannelExist
	}
	s.Channels[ch.ChannelID] = ch
	return nil
}
