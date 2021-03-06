package channel

import (
	"fmt"
	"log"

	model "github.com/socketChat/models"
)

var (
	ErrChannelExist    = fmt.Errorf("Channel exists")
	ErrChannelNotExist = fmt.Errorf("Channel does not exist")
)

type ChanMgr struct {
	addr     string
	channels map[string]*Channel
}

func NewChanMgr(addr string) *ChanMgr {
	return &ChanMgr{
		addr:     addr,
		channels: map[string]*Channel{},
	}
}

func (s *ChanMgr) NewChannel(channelName string) error {
	if _, exist := s.channels[channelName]; exist {
		return ErrChannelExist
	}
	ch := NewChannel(channelName)
	s.channels[channelName] = ch
	go ch.Serve()
	return nil
}

func (s *ChanMgr) GetChannel(channelID string) (*Channel, error) {
	ch, exist := s.channels[channelID]
	if !exist {
		return nil, ErrChannelNotExist
	}
	return ch, nil
}

func (s *ChanMgr) LeaveAllChannels(userID model.ID) error {
	for _, ch := range s.channels {
		if err := ch.Leave(userID); err != nil && err != ErrOpenIDNotExist {
			log.Println("Failed to leave channel", err.Error())
		}
	}
	return nil
}
