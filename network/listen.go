package network

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

// Listen ...
func Listen(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		// handle error
		fmt.Println("oh no!", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("oh no!", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// handle error
			fmt.Println("oh no!", err)
			return
		}

		fmt.Println(conn.RemoteAddr(), "said :", message)
	}
}
