package network

import (
	"fmt"
	"net"
	"strconv"
)

// Dial ...
func Dial(peer Peer) {
	conn, err := net.Dial("tcp", peer.address+":"+strconv.Itoa(peer.port))
	if err != nil {
		// handle error
		fmt.Println("error", err)
		return
	}

	for msg := range peer.Send {
		conn.Write([]byte(msg))
	}
}
