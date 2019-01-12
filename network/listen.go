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
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("oh no!", err)
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
		}

		fmt.Println(conn.RemoteAddr(), "said :", message)
	}
}
