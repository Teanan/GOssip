package chat

import (
	"fmt"

	"github.com/teanan/GOssip/network"
)

type MessageReceiver struct {
}

func (receiver *MessageReceiver) Receive(message network.Message, from network.Peer) {
	switch message.Kind {
	case "SAY":
		handleSay(message.Data, from)
	default:
		fmt.Println("Unknown message kind :", message)
	}
}

func handleSay(data string, from network.Peer) {
	fmt.Println("[", from, "] ", data)
}
