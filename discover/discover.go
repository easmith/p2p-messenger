package discover

import (
	"bufio"
	"github.com/easmith/p2p-messenger/proto"
	"log"
	"net"
	"os"
)

func StartDiscover(p *proto.Proto) {

	file, err := os.Open("./peers.txt")
	if err != nil {
		log.Printf("DISCOVER: Open peers.txt error: %s", err)
		return
	}

	var lastPeers []string

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
		log.Printf("Dial ERROR: " + err.Error())
		return
	}

	defer conn.Close()

	peer := handShake(p, conn)

	if peer == nil {
		log.Printf("Fail on handshake")
		return
	}

	p.PeerListener(peer)

	p.UnregisterPeer(peer)

	// TODO: ping-pong
}

func handShake(p *proto.Proto, conn net.Conn) *proto.Peer {

	p.SendName(conn)

	message, err := proto.ReadMessage(bufio.NewReader(conn))
	if err != nil {
		log.Printf("Error on read Message: %s", err)
		return nil
	}

	peer := proto.CreatePeer(message, conn)
	if peer != nil {
		p.RegisterPeer(peer)
	}

	// TODO: request peers

	return peer
}
