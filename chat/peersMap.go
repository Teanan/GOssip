package chat

import (
	"strconv"
	"strings"

	"github.com/teanan/GOssip-TP/network"
)

type peersMap struct {
	peers         map[string]network.Peer
	localUsername string
}

func (pmap *peersMap) Get(addr string) network.Peer {
	return pmap.peers[addr]
}

func (pmap *peersMap) Set(addr string, peer network.Peer) {
	// Question 1
}

func (pmap *peersMap) Find(address string) (bool, network.Peer) {
	// Question 1
	return false, network.Peer{}
}

func (pmap *peersMap) FindByName(name string) (bool, network.Peer) {
	// Question 3
	return false, network.Peer{}
}

func (pmap *peersMap) SendToAll(msg network.Message) {
	// Question 1
}

func (pmap *peersMap) SendTo(peer network.Peer, msg network.Message) {
	peer.Send <- msg
}

func (pmap *peersMap) SetNewPeersList(newList map[string]string, onPeerConnected func(network.Peer), onPeerDisconnected func(network.Peer)) {
	// remove peers that are no longer present
	for addr := range pmap.peers {
		_, found := newList[addr]
		if !found {
			onPeerDisconnected(pmap.peers[addr])
			delete(pmap.peers, addr)
		}
	}

	// add new peers
	for addr := range newList {
		_, found := pmap.peers[addr]
		if found {
			continue
		}

		port, _ := strconv.Atoi(strings.SplitN(addr, ":", 2)[1])
		peer := network.CreatePeer(
			strings.SplitN(addr, ":", 2)[0],
			port,
		)
		pmap.peers[addr] = peer
		onPeerConnected(peer)
	}
}

func (pmap *peersMap) SetLocalUsername(localUsername string) {
	pmap.localUsername = localUsername
}

func (pmap *peersMap) GetLocalUsername() string {
	return pmap.localUsername
}

func NewPeersMap() *peersMap {
	return &peersMap{
		peers: make(map[string]network.Peer),
	}
}
