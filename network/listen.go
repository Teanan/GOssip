package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type MessageReceiver interface {
	Receive(Message, Peer)
	HandleHello(data string, from Peer)
}

// Listen ...
func Listen(port int, peers PeersMap, messageReceiver MessageReceiver) {
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
		go handleConnection(conn, peers, messageReceiver)
	}
}

func handleConnection(conn net.Conn, peers PeersMap, messageReceiver MessageReceiver) {
	var remotePeerAddress = ""
	for {
		message, err := GetNextMessage(conn)
		if err != nil {
			fmt.Println("Failed to read message from peer", err)
			return
		}

		if message.Kind == "HELLO" {
			handleHello(message.Data, conn, peers, &remotePeerAddress)
			messageReceiver.HandleHello(message.Data, peers.Get(remotePeerAddress))
		} else {
			messageReceiver.Receive(message, peers.Get(remotePeerAddress))
		}
	}
}

func handleHello(data string, conn net.Conn, peers PeersMap, remotePeerAddress *string) {
	port, err := strconv.Atoi(strings.TrimSpace(data))

	if err != nil {
		fmt.Println("Invalid HELLO message ", err)
		return
	}

	addr := strings.Split(conn.RemoteAddr().String(), ":")[0] + ":" + strconv.Itoa(port)

	found, p := peers.Find(addr)

	if !found {
		fmt.Println("Unknown peer", addr)
		return
	}

	*remotePeerAddress = p.FullAddress()

	fmt.Println("Identified", conn.RemoteAddr(), "as", p)
}
