package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	connectedToDirectory = false
	peers                map[string]Peer
	peersMapChannel      chan map[string]Peer
	chatPort             int
)

// ConnectToDirectory ...
func ConnectToDirectory(directoryServer string, directoryPort int, localChatPort int, peersMap chan map[string]Peer) {
	peersMapChannel = peersMap
	peers = make(map[string]Peer)
	chatPort = localChatPort

	conn, err := net.Dial("tcp", directoryServer+":"+strconv.Itoa(directoryPort))
	if err != nil {
		fmt.Println("Cannot connect to directory", err)
		return
	}

	go listenFromDirectory(conn)
	connectedToDirectory = true
	send(conn, "HELLO", strconv.Itoa(chatPort))

	for {
		if !connectedToDirectory {
			conn, err := net.Dial("tcp", directoryServer+":"+strconv.Itoa(directoryPort))
			if err != nil {
				fmt.Println("Cannot connect to directory", err)
			} else {
				connectedToDirectory = true
				send(conn, "HELLO", strconv.Itoa(chatPort))
				go listenFromDirectory(conn)
			}

		}

		time.Sleep(2 * time.Second)
	}
}

func listenFromDirectory(conn net.Conn) {
	for {
		message, err := GetNextMessage(conn)
		if err != nil {
			fmt.Println("Lost connection to directory ", err)
			connectedToDirectory = false
			return
		}

		fmt.Println(conn.RemoteAddr(), "said :", message)

		switch message.Kind {
		case "PEERS":
			handlePeers(message.Data)

		case "NAME":
			handleName(message.Data)

		default:
			fmt.Println("Unknown message kind :", message)
		}
	}
}

func handlePeers(sList string) {
	list := strings.Split(sList, " ")

	newPeersList := make(map[string]string)

	// convert peers list to a map to simplify search by address
	// (and we only keep valid addresses)
	for _, addr := range list {
		if len(strings.SplitN(addr, ":", 2)) != 2 {
			continue
		}
		newPeersList[addr] = addr
	}

	// remove peers that are no longer present
	for addr := range peers {
		_, found := newPeersList[addr]
		if !found {
			fmt.Println(peers[addr], "left the chat")
			delete(newPeersList, addr)
		}
	}

	// add new peers
	for _, addr := range newPeersList {
		_, found := peers[addr]
		if found {
			continue
		}

		port, _ := strconv.Atoi(strings.SplitN(addr, ":", 2)[1])
		peer := Peer{
			address: strings.SplitN(addr, ":", 2)[0],
			port:    port,
			Send:    make(chan Message),
		}
		peers[addr] = peer
		fmt.Println(peers[addr], "joined the chat")
		go Dial(peer, chatPort)
	}

	peersMapChannel <- peers
}

func handleName(data string) {
	list := strings.SplitN(data, " ", 2)
	if len(list) < 2 {
		fmt.Println("Invalid NAME message", data)
		return
	}

	addr, newName := strings.TrimSpace(list[0]), strings.TrimSpace(list[1])

	peer := peers[addr]
	peer.name = newName

	fmt.Println(peers[addr], "is now", peer)

	peers[addr] = peer
}

func send(conn net.Conn, msgType string, data string) (int, error) {
	return conn.Write([]byte(msgType + " " + data + "\n"))
}
