package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/teanan/GOssip/chat"
	"github.com/teanan/GOssip/network"
)

func main() {
	fmt.Println("== GOssip ==")

	rand.Seed(time.Now().UnixNano())
	port := 9000 + rand.Intn(1000)

	directoryPort := 8080
	directoryServer := "127.0.0.1"

	fmt.Println("Listening on port", port)

	peersMap := chat.NewPeersMap()
	commandProcessor := chat.NewCommandProcessor(peersMap)

	peersMapChannel := make(chan map[string]network.Peer)

	go network.Listen(port, peersMap, &chat.MessageReceiver{})
	go network.ConnectToDirectory(directoryServer, directoryPort, port, peersMapChannel)

	stdin := make(chan string)
	go readStdin(stdin)

	for {

		select {

		case text, ok := <-stdin: // New command from stdin

			if !ok {
				return
			}

			commandProcessor.Process(text)

		case newMap := <-peersMapChannel: // New peers list from discovery server
			peersMap.SetAll(newMap)

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
