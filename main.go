package main

import (
	"log"
	"net"
	"time"

	"github.com/vincepr/go_distributed_cache/cache"
)


func main(){
	options := ServerOpts{
		ListenAddr: ":5555",
		IsLeader: true,
	}

	// testing client only
	go func (){
		time.Sleep(time.Second)
		con, err := net.Dial("tcp", ":5555")
		if err != nil{
			log.Fatal(err)
		}
		con.Write([]byte("SET Foo Bar 2500"))
	}()


	// our hard coded leader server
	server := NewServer(options, cache.NewCache())
	server.Start()
}