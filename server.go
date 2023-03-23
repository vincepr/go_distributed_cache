package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/vincepr/go_distributed_cache/cache"
)

type ServerOpts struct{
	ListenAddr 	string
	IsLeader	bool
}

type Server struct{
	ServerOpts

	cache cache.Cacher
}

func NewServer(options ServerOpts, c cache.Cacher) *Server{
	return &Server{
		ServerOpts: options,
		cache: c,
	}
}

func (s *Server) Start() error{
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil{
		return err
	}
	log.Println("server listening on port:", s.ListenAddr)

	for {
		con, err := listener.Accept()
		if err != nil{
			log.Println("connection error:", err)
			continue
		}
		go s.handleConnection(con)	// pass each connection to a new goroutine
	}
}

func (s *Server) handleConnection(con net.Conn){
	defer func(){
		con.Close()
	}()

	buf := make([]byte, 2048)		// 2048 max size of our buffer
	for {
		n, err := con.Read(buf)
		if err != nil {
			log.Println("handleConnection read-error:", err)
			break					// after a connection error we force the connection close
		}
		data := buf[:n]							//:todo remove
		fmt.Println("incoming: ",string(data))	//:todo remove
		go s.handleCommand(con, buf[:n])
	}
}

func (s *Server) handleCommand(con net.Conn, byteCmd []byte){
	msg, err := parseCommand(byteCmd)
	if err != nil{
		log.Println(err)	//:todo when responding to client implemented
		//:todo respond to client
		return
	}
	switch msg.Cmd{
	case CMDSet:
		if err := s.handleSetCommand(con, msg);err != nil{
			//:todo respond to client
			return
		}
	}

	

}

func (s *Server) handleSetCommand(con net.Conn,msg *Message) error{
	fmt.Println("handling the following SET command:", msg)	//:todo remove debug/crafting only

	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err!= nil{
		return err
	}

	go s.sendToFollowers(context.TODO())


	return nil
}

func (s *Server) sendToFollowers(ctx context.Context) error{

	return nil
}
