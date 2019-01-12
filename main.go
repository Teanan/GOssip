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
	port := 8000 + rand.Intn(2000)

	fmt.Println("Listening on port", port)

	peersMapChannel := make(chan map[string]network.Peer)
	peersMap := make(map[string]network.Peer)

	go network.Listen(port)
	go network.DiscoverLoop(port, peersMapChannel)

	reader := bufio.NewReader(os.Stdin)

	for {

		select {
		case newMap := <-peersMapChannel:
			peersMap = newMap
		}

		text, _ := reader.ReadString('\n')

		for _, peer := range peersMap {
			peer.Send <- text
		}
	}
}
