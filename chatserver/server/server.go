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
	chanMgr      *channel.ChanMgr
	listener     *net.TCPListener
	connID       int32
	shutdownChan chan interface{}
}

func NewServer(host, port string) (*Server, error) {
	addr := host + ":" + port
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	listener, err := net.ListenTCP("tcp", tcpAddr)
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
		chanMgr:      chanMgr,
		listener:     listener,
		shutdownChan: make(chan interface{}),
	}
	return server, nil
}

func (s *Server) Serve() {
	defer log.Println("Server shutdown")
	// spawn a goroutine to handle incoming connector
	connChan := make(chan net.Conn)
	go func() {
		defer s.ShutDown()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Println(err.Error())
				return
			}
			connChan <- conn
		}
	}()
	// Process connection and new client
	for {
		select {
		case <-s.shutdownChan:
			s.listener.Close()
			return
		case conn := <-connChan:
			clientconn := clientconn.NewClient(conn, s.connID, s.chanMgr)
			go clientconn.Listen()
			s.connID++
		}
	}
}

func (s *Server) ShutDown() {
	s.shutdownChan <- struct{}{}
}
