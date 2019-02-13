package server

import (
	"context"
	"log"
	"net"
	"net/http"

	channel "github.com/socketChat/chatserver/channel"
	"github.com/socketChat/chatserver/clientconn"
	"golang.org/x/net/websocket"
)

const (
	defaultChannel = "happy-pig-year"
)

type Server interface {
	Serve()
	ShutDown()
}

// SocketServer defines socket based chatroom server
type SocketServer struct {
	chanMgr      *channel.ChanMgr
	listener     *net.TCPListener
	connID       int32
	shutdownChan chan interface{}
}

func NewSocketServer(host, port string) (Server, error) {
	addr := host + ":" + port
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	chanMgr := channel.NewChanMgr(addr)
	// Init a chat room and run
	if err := chanMgr.NewChannel(defaultChannel); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	server := &SocketServer{
		chanMgr:      chanMgr,
		listener:     listener,
		shutdownChan: make(chan interface{}),
	}
	return server, nil
}

func (s *SocketServer) Serve() {
	defer log.Println("Server shutdown")
	// spawn a goroutine to handle incoming connector
	go func() {
		defer s.ShutDown()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Println(err.Error())
				return
			}
			clientconn := clientconn.NewClient(conn, s.connID, s.chanMgr)
			go clientconn.Listen()
			s.connID++
		}
	}()
	// Process connection and new client
	<-s.shutdownChan
	s.listener.Close()
}

func (s *SocketServer) ShutDown() {
	s.shutdownChan <- struct{}{}
}

// WebSocketServer defines websocket based chatroom server
type WebSocketServer struct {
	chanMgr      *channel.ChanMgr
	connID       int32
	httpserver   *http.Server
	shutdownChan chan interface{}
}

func (wss *WebSocketServer) handler(conn *websocket.Conn) {
	clientconn := clientconn.NewClient(conn, wss.connID, wss.chanMgr)
	wss.connID++
	clientconn.Listen()
}

func NewWebSocketServer(port string) (Server, error) {
	chanMgr := channel.NewChanMgr("")
	// Init a chat room and run
	if err := chanMgr.NewChannel(defaultChannel); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	server := &WebSocketServer{
		chanMgr:      chanMgr,
		httpserver:   &http.Server{Addr: ":" + port},
		shutdownChan: make(chan interface{}),
	}
	return server, nil
}

func (wss *WebSocketServer) Serve() {
	defer log.Println("Server shutdown")
	go func() {
		http.Handle("/chat", websocket.Handler(wss.handler))
		if err := wss.httpserver.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	<-wss.shutdownChan
	wss.httpserver.Shutdown(context.Background())
}

func (wss *WebSocketServer) ShutDown() {
	wss.shutdownChan <- struct{}{}
}
