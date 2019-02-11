package server

import (
	"log"
	"net"

	channel "github.com/socketChat/chatserver/channel"
	"github.com/socketChat/chatserver/clientconn"
)

const (
	defaultChannel = "happy-pig-year"
)

type Server struct {
	chanMgr  *channel.ChanMgr
	listener net.Listener
	connID   int32
}

func NewServer(host, port string) (*Server, error) {
	addr := host + ":" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	chanMgr := channel.NewChanMgr(addr)
	// Init a chat room and run
	ch := channel.NewChannel(defaultChannel)
	go ch.Serve()
	if err := chanMgr.AddChannel(ch); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	server := &Server{
		chanMgr:  chanMgr,
		listener: listener,
	}
	return server, nil
}

func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err.Error())
			break
		}
		clientconn := clientconn.NewClient(conn, s.connID, s.chanMgr)
		go clientconn.Listen()
		s.connID++
	}
}
