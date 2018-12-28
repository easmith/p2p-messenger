package discover

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"github.com/easmith/p2p-messanger/proto"
	"log"
	"net"
	"os"
	"time"
)

func StartDiscover(p *proto.Proto) {

	file, err := os.Open("./peers.txt")
	if err != nil {
		log.Printf("DISCOVER: Open peers.txt error: %s", err)
		return
	}

	lastPeers := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastPeers = append(lastPeers, scanner.Text())
	}

	log.Printf("DISCOVER: Start peer discovering. Last seen peers: %v", len(lastPeers))
	for _, peerAddress := range lastPeers {
		go checkPeer(p, peerAddress)
	}
}

func checkPeer(p *proto.Proto, peerAddress string) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		log.Printf("DISCOVER: Dial Error: " + err.Error())
		return
	}

	defer conn.Close()

	peerName := handShake(p, conn)

	log.Printf("DISCOVER: Peer %s is ok: %s", peerAddress, peerName)

	hexPubKey, err := hex.DecodeString(peerName.PubKey)
	if err != nil {
		log.Printf("DISCOVERY: hex decode error: %s", err)
	}
	peer := &proto.Peer{
		PubKey:    hexPubKey,
		Addr:      "addr",
		Conn:      &conn,
		Name:      peerName.Name,
		FirstSeen: time.Now().String(),
		LastSeen:  time.Now().String(),
		Peers:     proto.NewPeers(),
	}
	p.Peers.Put(peer)

	proto.ConnListener(conn, p)

	// TODO: ping-pong
	// TODO: request peers
	// TODO: listenPeer
}

func handShake(p *proto.Proto, conn net.Conn) *proto.PeerName {

	p.SendName(conn)

	message, err := proto.ReadMessage(bufio.NewReader(conn))
	if err != nil {
		log.Printf("DISCOVER: Error on read Message: %s", err)
		return nil
	}

	log.Printf("DISCOVER: Peer Message: %s %s", message.Cmd, message.Content)

	peerName := proto.PeerName{}
	if string(message.Cmd) == "NAME" {
		err := json.Unmarshal(message.Content, &peerName)
		if err != nil {
			log.Printf("DISCOVER: error: %v", err)
			return nil
		}
	}

	return &peerName
}
