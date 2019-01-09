package proto

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"time"
)

func CreatePeer(message *Message, conn net.Conn) *Peer {
	if string(message.Cmd) != "NAME" {
		return nil
	}

	peerName := &PeerName{}
	err := json.Unmarshal(message.Content, peerName)
	if err != nil {
		log.Printf("handShake unmarshall ERROR: %v", err)
		return nil
	}

	rawPubKey, err := hex.DecodeString(peerName.PubKey)
	if err != nil {
		log.Printf("handShake hex decode error: %s", err)
		return nil
	}

	// TODO: проверить подпись

	peer := &Peer{
		PubKey:    rawPubKey,
		Addr:      conn.RemoteAddr().String(),
		Conn:      &conn,
		Name:      peerName.Name,
		FirstSeen: time.Now().String(),
		LastSeen:  time.Now().String(),
		Peers:     NewPeers(),
	}

	return peer
}
