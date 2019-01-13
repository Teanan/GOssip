package network

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Dial ...
func Dial(peer Peer, localChatPort int) {
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", peer.address+":"+strconv.Itoa(peer.port))
	if err != nil {
		fmt.Println("Failed to connect to peer", err)
		return
	}

	Message{
		"HELLO",
		strconv.Itoa(localChatPort),
	}.Send(conn)

	for msg := range peer.Send {
		msg.Send(conn)
	}
}
