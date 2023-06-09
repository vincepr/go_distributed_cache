package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
**	Defining the syntax of our cache
 */
type Command string

// all supported Key Keywords
const(
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)


type Message struct{
	Cmd 	Command
	Key 	[]byte
	Value 	[]byte
	TTL 	time.Duration
	
}

func (m *Message) ToBytes() []byte{
	switch m.Cmd{
	case CMDSet:
		cmd := []byte(fmt.Sprintf("%s %s %s %d", m.Cmd, m.Key, m.Value, m.TTL))
		return []byte(cmd)
	case CMDGet:
		cmd := []byte(fmt.Sprintf("%s %s", m.Cmd, m.Key))
		return []byte(cmd)
	default:
		panic("unknown command :todo")
	}

}

// parses our recieved bytes for our custom command syntax
func parseCommand(byteCmd []byte) (*Message, error){
	str := string(byteCmd)
	parts := strings.Split(str, " ")
	if len(parts) <2{			//:todo this might have to change for GET etc.
		return nil, errors.New("invalid protocol syntax. 0 Arguments.")
	}
	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte( parts[1] ),
	}
	//* SET - LOGIC*/
	if msg.Cmd == CMDSet{
		fmt.Println((parts))
		if len(parts) != 4{
			return nil, errors.New("invalid protocol syntax. SET must have 3 Arguments")
		}

		msg.Value = []byte(parts[2])
		ttl, err := strconv.Atoi((parts[3]) )
		if err != nil{
			return nil, errors.New("invalid protocol syntax. Cant parse Time-To-Live.")
		}
		msg.TTL = time.Duration(ttl)
	}
	return msg, nil
}