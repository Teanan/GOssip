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

func main() {
	fmt.Println("== GOssip ==")

	rand.Seed(time.Now().UnixNano())
	port := 9000 + rand.Intn(1000)

	directoryPort := 8080
	directoryServer := "127.0.0.1"

	fmt.Println("Listening on port", port)

	browserPort := 13000 + rand.Intn(1000)
	webpage, err := browser.Connect("localhost", browserPort)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Successfuly connected to browser webpage")

	peersMapChannel := make(chan map[string]network.Peer)
	usernameChannel := make(chan string)
	messageOutputChannel := make(chan string)

	peersMap := chat.NewPeersMap()
	commandProcessor := chat.NewCommandProcessor(peersMap, messageOutputChannel)

	go network.Listen(port, peersMap, chat.NewMessageReceiver(peersMap, messageOutputChannel))
	go network.ConnectToDirectory(directoryServer, directoryPort, port, peersMapChannel, usernameChannel)

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

		case newMap := <-peersMapChannel: // New peers list from discovery server
			peersMap.SetAll(newMap)

		case name := <-usernameChannel: // Assigned username from discovery server
			peersMap.SetLocalUsername(name)

		case messageFromBrowser := <-webpage.ReceiveChannel():
			messageOutputChannel <- "[YOU] " + messageFromBrowser
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
