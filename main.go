package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/teanan/GOssip-TP/chat"
	"github.com/teanan/GOssip-TP/network"
)

var (
	chatPort        int
	directoryPort   = 8080
	directoryServer = "127.0.0.1"

	messageOutputChannel = make(chan string, 5)
)

func main() {
	fmt.Println("== GOssip ==")

	rand.Seed(time.Now().UnixNano())
	chatPort = 9000 + rand.Intn(1000)

	if len(os.Args) > 1 {
		directoryServer = os.Args[1]
	}

	fmt.Println("Listening on port", chatPort)

	/* Question 4
	browserPort := 13000 + rand.Intn(1000)
	webpage, err := browser.Connect("localhost", browserPort)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfuly connected to browser webpage")
	*/

	peersListChannel := make(chan map[string]string, 5)
	usernameChannel := make(chan string, 5)

	peersMap := chat.NewPeersMap()
	commandProcessor := chat.NewCommandProcessor(peersMap, messageOutputChannel)

	go network.Listen(chatPort, peersMap, chat.NewMessageReceiver(peersMap, messageOutputChannel))
	go network.ConnectToDirectory(directoryServer, directoryPort, chatPort, peersListChannel, usernameChannel)

	stdin := make(chan string)
	go readStdin(stdin)

	/* Question 4
	loop:
	*/
	for {

		select {

		case text, ok := <-stdin: // New command from stdin

			if !ok {
				return
			}

			commandProcessor.Process(text)

		case newList := <-peersListChannel: // New peers list from discovery server
			peersMap.SetNewPeersList(newList, onPeerConnected, onPeerDisconnected)

		case name := <-usernameChannel: // Assigned username from discovery server
			peersMap.SetLocalUsername(name)

		case message := <-messageOutputChannel:
			fmt.Println(message)

			/* Question 4
			case <-webpage.Disconnected:
				fmt.Println("Browser webpage has disconnected")
				break loop
			*/
		}
	}
}

func onPeerConnected(peer network.Peer) {
	// Question 2
}

func onPeerDisconnected(peer network.Peer) {
	// Question 2
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
