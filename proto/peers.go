package proto

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
	"log"
	"net"
	"sync"
)

type Peer struct {
	PubKey    ed25519.PublicKey
	Addr      string
	Conn      *net.Conn
	Name      string
	FirstSeen string
	LastSeen  string
	Peers     *Peers
}

func (p Peer) String() string {
	peerName, err := json.Marshal(PeerName{
		Name:   p.Name,
		PubKey: hex.EncodeToString(p.PubKey),
	})

	if err != nil {
		panic(err)
	}

	return string(peerName)
}

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

func (p Peers) Remove(peer *Peer) (found bool) {
	p.RLock()
	defer p.RUnlock()

	delete(p.peers, string(peer.PubKey))
	return
}

func (p Peers) PeerList() *PeerList {

	peerList := &PeerList{
		WsCmd: WsCmd{
			Cmd: "PEERS",
		},
		Peers: []PeerName{},
	}

	log.Printf("total peers: %v", len(p.peers))

	for _, el := range p.peers {
		peerList.Peers = append(peerList.Peers, PeerName{
			Name:   el.Name,
			PubKey: hex.EncodeToString(el.PubKey),
		})
	}

	return peerList
}
