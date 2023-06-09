package main

import (
	"flag"

	"github.com/vincepr/go_distributed_cache/cache"
)


func main(){
	// // TestingONLY: checking if the set command gets forwarded from the leader to the followers
	// connetion, err := net.Dial("tcp", ":6666")
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// _, err = connetion.Write([]byte("SET Foo Bar 99999999"))
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// return

	var(
		listenAddr = flag.String("listen", "5555", "listen port of the service")
		leaderAddr = flag.String("leaderaddr", "", "listen adress of the leader")
	)
	
	flag.Parse()

	options := ServerOpts{
		ListenAddr: ":" + *listenAddr,
		IsLeader: *leaderAddr=="",
		LeaderAddr: ":"+ *leaderAddr,
	}

	// // testing client only
	// go func (){
	// 	time.Sleep(time.Second)
	// 	con, err := net.Dial("tcp", ":5555")
	// 	if err != nil{
	// 		log.Fatal(err)
	// 	}
	// 	con.Write([]byte("SET Foo Bar 2500000000"))

	// 	time.Sleep(time.Second)
	// 	con.Write([]byte("GET Foo"))
	// 	buf := make([]byte, 1000)
	// 	n, _ := con.Read(buf)
	// 	log.Println(string(buf[:n]))
	// }()


	// our hard coded leader server
	server := NewServer(options, cache.NewCache())
	server.Start()
	
}