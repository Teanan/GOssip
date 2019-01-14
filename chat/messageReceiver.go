package chat

import (
	"fmt"
	"strings"

	"github.com/teanan/GOssip/network"
)

type MessageReceiver struct {
	peers         *peersMap
	messageOutput chan<- string
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

func (receiver *MessageReceiver) HandleHello(data string, from network.Peer) {
	receiver.peers.SendTo(from, network.Message{
		Kind: "NAME",
		Data: receiver.peers.GetLocalUsername(),
	})
}

func (receiver *MessageReceiver) handleSay(data string, from network.Peer) {
	receiver.messageOutput <- fmt.Sprint("[", from, "] ", data)
}

func (receiver *MessageReceiver) handleSayTo(data string, from network.Peer) {
	// question 5
}

func (receiver *MessageReceiver) handleName(data string, from network.Peer) {
	if strings.ContainsAny(data, "\t\r\n ") {
		return
	}

	if data == from.String() {
		return
	}

	if found, _ := receiver.peers.FindByName(data); found || receiver.peers.GetLocalUsername() == data {
		receiver.messageOutput <- fmt.Sprint(from.String(), "tried to use an already taken username")
		return
	}

	receiver.messageOutput <- fmt.Sprint(from.String(), " is now known as ", data)
	from.SetName(data)
	receiver.peers.Set(from.FullAddress(), from)
}

func NewMessageReceiver(peers *peersMap, messageOutput chan<- string) *MessageReceiver {
	return &MessageReceiver{
		peers:         peers,
		messageOutput: messageOutput,
	}
}
