package chat

import "github.com/teanan/GOssip/network"

var (
	Peers peersMap
)

type peersMap struct {
	peers map[string]network.Peer
}

func (pmap *peersMap) Get(addr string) network.Peer {
	return pmap.peers[addr]
}

func (pmap *peersMap) Find(address string) (bool, network.Peer) {
	for addr, peer := range pmap.peers {
		if addr == address {
			return true, peer
		}
	}
	return false, network.Peer{}
}

func (pmap *peersMap) SendToAll(text string) {
	for _, peer := range pmap.peers {
		peer.Send <- network.Message{
			Kind: "SAY",
			Data: text,
		}
	}
}

func (pmap *peersMap) SetAll(newMap map[string]network.Peer) {
	pmap.peers = newMap
}

func NewPeersMap() *peersMap {
	return &peersMap{
		peers: make(map[string]network.Peer),
	}
}
