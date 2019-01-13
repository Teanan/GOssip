package chat

import "github.com/teanan/GOssip/network"

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

func (pmap *peersMap) FindByName(name string) (bool, network.Peer) {
	for _, peer := range pmap.peers {
		if peer.Name() == name {
			return true, peer
		}
	}
	return false, network.Peer{}
}

func (pmap *peersMap) SendToAll(msg network.Message) {
	for _, peer := range pmap.peers {
		peer.Send <- msg
	}
}

func (pmap *peersMap) SendTo(peer network.Peer, msg network.Message) {
	peer.Send <- msg
}

func (pmap *peersMap) SetAll(newMap map[string]network.Peer) {
	pmap.peers = newMap
}

func NewPeersMap() *peersMap {
	return &peersMap{
		peers: make(map[string]network.Peer),
	}
}
