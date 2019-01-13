package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/teanan/GOssip/browser"
	"github.com/teanan/GOssip/chat"
	"github.com/teanan/GOssip/network"
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

	fmt.Println("Listening on port", chatPort)

	browserPort := 13000 + rand.Intn(1000)
	webpage, err := browser.Connect("localhost", browserPort)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfuly connected to browser webpage")

	peersListChannel := make(chan map[string]string, 5)
	usernameChannel := make(chan string, 5)

	peersMap := chat.NewPeersMap()
	commandProcessor := chat.NewCommandProcessor(peersMap, messageOutputChannel)

	go network.Listen(chatPort, peersMap, chat.NewMessageReceiver(peersMap, messageOutputChannel))
	go network.ConnectToDirectory(directoryServer, directoryPort, chatPort, peersListChannel, usernameChannel)

	stdin := make(chan string)
	go readStdin(stdin)

loop:
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

		case messageFromBrowser := <-webpage.ReceiveChannel():
			commandProcessor.Process(messageFromBrowser)

		case message := <-messageOutputChannel:
			fmt.Println(message)
			webpage.SendChannel() <- message

		case <-webpage.Disconnected:
			fmt.Println("Browser webpage has disconnected")
			break loop
		}
	}
}

func onPeerConnected(peer network.Peer) {
	messageOutputChannel <- fmt.Sprint(peer, " joined the chat")
	go network.Dial(peer, chatPort)
}

func onPeerDisconnected(peer network.Peer) {
	messageOutputChannel <- fmt.Sprint(peer, " left the chat")
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
