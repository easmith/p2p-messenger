package proto

import (
	"sync"
)

type Peers struct {
	sync.RWMutex
	peers map[string]*Peer
}

func NewPeers() *Peers {
	return &Peers{
		peers: make(map[string]*Peer),
	}
}

func (p Peers) Put(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	p.peers[string(peer.PubKey)] = peer
}

func (p Peers) Get(key string) (peer *Peer, found bool) {
	p.RLock()
	defer p.RUnlock()

	peer, found = p.peers[key]
	return
}

func (p Peers) List() *map[string]*Peer {
	return &p.peers
}
