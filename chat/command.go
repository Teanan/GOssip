package chat

import (
	"fmt"
	"strings"

	"github.com/teanan/GOssip/network"
)

type commandProcessor struct {
	peers         *peersMap
	messageOutput chan<- string
}

func (processor *commandProcessor) Process(command string) {

	if strings.HasPrefix(command, "/") {
		commandName := ""

		// question 5

		switch commandName {
		default:
			fmt.Print("Unknown command", commandName)
		}
	} else {
		processor.messageOutput <- "[" + processor.peers.GetLocalUsername() + "] " + command
		processor.peers.SendToAll(network.Message{
			Kind: "SAY",
			Data: command,
		})
	}
}

func (processor *commandProcessor) say(commandParams string) {

}

func NewCommandProcessor(peers *peersMap, messageOutput chan<- string) *commandProcessor {
	return &commandProcessor{
		peers:         peers,
		messageOutput: messageOutput,
	}
}
