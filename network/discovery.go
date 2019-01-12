package network

import (
	"fmt"
	"strconv"
	"time"

	"github.com/schollz/peerdiscovery"
)

// Peer ...
type Peer struct {
	address string
	port    int
	Send    chan string
}

func discover(port int) map[string]Peer {
	discoveries, _ := peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     50,
		Delay:     1,
		TimeLimit: 2,
		AllowSelf: true,
		Payload:   []byte(strconv.Itoa(port)),
	})

	peers := make(map[string]Peer)

	for _, d := range discoveries {
		pport, _ := strconv.Atoi(string(d.Payload[:4]))
		peers[d.Address] = Peer{
			d.Address,
			pport,
			make(chan string),
		}
	}

	return peers
}

// DiscoverLoop runs peers discovery every 2 seconds (Function to be ran in parallel)
func DiscoverLoop(port int, peersMap chan map[string]Peer) {
	peers := make(map[string]Peer)
	for {
		newPeers := discover(port)

		for key := range peers {
			_, found := newPeers[key]
			if found {
				delete(newPeers, key)
			}
		}

		for key, peer := range newPeers {
			fmt.Println(key+":"+strconv.Itoa(peer.port), "joined the chat")
			peers[key] = peer
			go Dial(peer)
		}

		peersMap <- peers

		time.Sleep(2 * time.Second)
	}
}
