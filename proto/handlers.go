package proto

import (
	"log"
	"time"
)

func (p Proto) onHand(peer *Peer, envelope *Envelope) {
	log.Printf("onHand")
	newPeer := NewPeer(*peer.Conn)

	err := newPeer.UpdatePeer(envelope)
	if err != nil {
		log.Printf("Update peer error: %s", err)
	} else {
		if peer != nil {
			p.UnregisterPeer(peer)
		}

		// TODO: не заменяется по ссылке (peer = newPeer), приходится копировать поля
		peer.Name = newPeer.Name
		peer.PubKey = newPeer.PubKey
		peer.SharedKey = newPeer.SharedKey
		peer.LastSeen = time.Now().String()

		p.RegisterPeer(peer)
	}
	p.SendName(peer)
	return
}

func (p Proto) onMess(peer *Peer, envelope *Envelope) {
	envelope.Content = Decrypt(envelope.Content, peer.SharedKey.Secret)
	p.Broker <- envelope
	return
}

func (p Proto) onList(peer *Peer, envelope *Envelope) {
	log.Printf("onList")
	return
}
