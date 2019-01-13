package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Listen ...
func Listen(port int, peerFinder func(string) (bool, Peer)) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to open listen socket", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection peer", err)
			return
		}
		go handleConnection(conn, peerFinder)
	}
}

func handleConnection(conn net.Conn, peerFinder func(string) (bool, Peer)) {
	var remotePeer Peer
	for {
		message, err := GetNextMessage(conn)
		if err != nil {
			fmt.Println("Failed to read message from peer", err)
			return
		}

		switch message.Kind {
		case "HELLO":
			handleHello(message.Data, conn, peerFinder, &remotePeer)
		case "SAY":
			handleSay(message.Data, remotePeer)
		default:
			fmt.Println("Unknown message kind :", message)
		}
	}
}

func handleHello(data string, conn net.Conn, peerFinder func(string) (bool, Peer), peer *Peer) {
	port, err := strconv.Atoi(strings.TrimSpace(data))

	if err != nil {
		fmt.Println("Invalid HELLO message ", err)
		return
	}

	addr := strings.Split(conn.RemoteAddr().String(), ":")[0] + ":" + strconv.Itoa(port)

	found, p := peerFinder(addr)

	if !found {
		fmt.Println("Unknown peer", addr)
		return
	}

	*peer = p
	fmt.Println("Identified", conn.RemoteAddr(), "as", peer)
}

func handleSay(data string, peer Peer) {
	fmt.Println("[", peer, "] ", data)
}
