package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/teanan/GOssip/network"
)

func main() {
	fmt.Println("== GOssip ==")

	rand.Seed(time.Now().UnixNano())
	port := 9000 + rand.Intn(1000)

	directoryPort := 8080
	directoryServer := "127.0.0.1"

	fmt.Println("Listening on port", port)

	peersMapChannel := make(chan map[string]network.Peer)
	peersMap := make(map[string]network.Peer)

	go network.Listen(port, func(address string) (bool, network.Peer) {
		for addr, peer := range peersMap {
			if addr == address {
				return true, peer
			}
		}
		return false, network.Peer{}
	})
	go network.ConnectToDirectory(directoryServer, directoryPort, port, peersMapChannel)

	stdin := make(chan string)
	go readStdin(stdin)

	for {

		select {

		case text, ok := <-stdin: // New message from stdin

			if !ok {
				return
			}

			for _, peer := range peersMap {
				peer.Send <- network.Message{
					"SAY",
					text,
				}
			}

		case newMap := <-peersMapChannel: // New peers list from discovery server
			peersMap = newMap

		}

	}
}

func readStdin(ch chan string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			close(ch)
			return
		}
		ch <- s
	}
}
