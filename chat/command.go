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
		processor.messageOutput <- "{CMD} " + command
		command = strings.TrimSpace(command)
		commandName := strings.SplitN(command, " ", 2)[0]
		commandParams := ""
		if len(strings.SplitN(command, " ", 2)) > 1 {
			commandParams = strings.TrimSpace(strings.SplitN(command, " ", 2)[1])
		}

		switch commandName {
		case "/say":
			processor.say(commandParams)
		case "/name":
			processor.name(commandParams)
		default:
			fmt.Print("Unknown command", commandName)
		}
	} else {
		processor.messageOutput <- "You said: " + command
		processor.peers.SendToAll(network.Message{
			Kind: "SAY",
			Data: command,
		})
	}
}

func (processor *commandProcessor) say(commandParams string) {
	if len(strings.SplitN(commandParams, " ", 2)) != 2 {
		processor.messageOutput <- fmt.Sprint("Usage: /say username message")
		return
	}

	userName, text := strings.SplitN(commandParams, " ", 2)[0], strings.SplitN(commandParams, " ", 2)[1]

	ok, peer := processor.peers.FindByName(userName)

	if !ok {
		processor.messageOutput <- fmt.Sprint("Cannot find user", userName)
		return
	}

	processor.peers.SendTo(peer, network.Message{
		Kind: "SAYTO",
		Data: text,
	})
}

func (processor *commandProcessor) name(commandParams string) {
	if commandParams == "" || len(strings.Split(commandParams, " ")) != 1 {
		processor.messageOutput <- "Usage: /name new_username"
		return
	}

	if found, _ := processor.peers.FindByName(commandParams); found {
		processor.messageOutput <- "Username already taken"
		return
	}

	processor.peers.SendToAll(network.Message{
		Kind: "NAME",
		Data: commandParams,
	})

	processor.peers.SetLocalUsername(commandParams)
	processor.messageOutput <- "Thou shall now be called " + commandParams
}

func NewCommandProcessor(peers *peersMap, messageOutput chan<- string) *commandProcessor {
	return &commandProcessor{
		peers:         peers,
		messageOutput: messageOutput,
	}
}
