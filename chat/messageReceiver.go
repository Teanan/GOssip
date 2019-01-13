package chat

import (
	"fmt"

	"github.com/teanan/GOssip/network"
)

type MessageReceiver struct {
	peers *peersMap
}

func (receiver *MessageReceiver) Receive(message network.Message, from network.Peer) {
	switch message.Kind {
	case "SAY":
		receiver.handleSay(message.Data, from)
	case "SAYTO":
		receiver.handleSayTo(message.Data, from)
	case "NAME":
		receiver.handleName(message.Data, from)
	default:
		fmt.Println("Unknown message kind :", message)
	}
}

func (receiver *MessageReceiver) handleSay(data string, from network.Peer) {
	fmt.Println("[", from, "] ", data)
}

func (receiver *MessageReceiver) handleSayTo(data string, from network.Peer) {
	fmt.Println("/", from, "/ ", data)
}

func (receiver *MessageReceiver) handleName(data string, from network.Peer) {
	if found, _ := receiver.peers.FindByName(data); found {
		fmt.Println(from.Name(), "tried to use an already taken username")
		return
	}

	fmt.Println(from.Name(), " is now known as ", data)
	from.SetName(data)
	receiver.peers.Set(from.FullAddress(), from)
}

func NewMessageReceiver(peers *peersMap) *MessageReceiver {
	return &MessageReceiver{
		peers: peers,
	}
}
