package network

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Dial ...
func Dial(peer Peer, localChatPort int) {
	conn, err := net.Dial("tcp", peer.address+":"+strconv.Itoa(peer.port))
	if err != nil {
		fmt.Println("Failed to connect to peer", err)
		return
	}

	fmt.Println("Connected to", peer.address+":"+strconv.Itoa(peer.port))

	Message{
		"HELLO",
		strconv.Itoa(localChatPort),
	}.Send(conn)

	for {
		select {
		case msg := <-peer.Send:
			msg.Send(conn)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
