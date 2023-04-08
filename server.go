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
	LeaderAddr	string
}

type Server struct{
	ServerOpts
	cache cache.Cacher						// our own cach we save/read data from
	followers map[net.Conn] struct{}		// stores all connected child-caches
}

func NewServer(options ServerOpts, c cache.Cacher) *Server{
	return &Server{
		ServerOpts: options,
		cache: c,
		// :todo only allocate this when we are the leader?
		followers: 	make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error{
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil{
		return err
	}
	log.Println("server listening on port:", s.ListenAddr)

	if !s.IsLeader{
		go func() {
			con, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Println("connected with leader:", s.LeaderAddr)
			if err != nil{
				log.Fatal(err)
			}
			s.handleConnection(con)
		}()
	}

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
	defer con.Close()
	fmt.Println("connection made with:", con.RemoteAddr())

	buf := make([]byte, 2048)		// 2048 max size of our buffer

	if s.IsLeader {
		s.followers[con] = struct{}{}
	}

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
		log.Println("Failed to parse command",err)	//:todo remove when responding to client implemented
		con.Write([]byte(err.Error()))
		return
	}

	log.Printf("hadnleCommand() recieved command %s", msg.Cmd)	//:todo remove this once everything is running
	// different supported methods -> handlers:

	switch msg.Cmd{
	case CMDSet:
		err = s.handleSetCommand(con, msg)
	case CMDGet:
		err = s.handleGetCommand(con, msg)
	}
	if err != nil{
		log.Println("Failed to handle command",err)	//:todo remove when responding to client implemented
		con.Write([]byte(err.Error()))
		return								//:todo remove this?
	}
	

}

func (s *Server) handleSetCommand(con net.Conn,msg *Message) error{
	fmt.Println("handling the following SET command:", msg)	//:todo remove debug/crafting only

	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err!= nil{
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)


	return nil
}

func (s *Server) handleGetCommand(con net.Conn,msg *Message) error{
	val, err := s.cache.Get(msg.Key)
	if err != nil{
		return err
	}
	if _, err := con.Write(val); err != nil{
		return err
	}
	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error{
	// :todo make chaning followers & reading (like here) mutex-protected. (does atomic map exist?)
	for con := range s.followers {
		_, err := con.Write(msg.ToBytes())
		if err != nil{
			log.Println("Error writing to follower", err)
			continue
		}
	}
	return nil
}
