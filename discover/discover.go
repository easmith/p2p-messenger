package discover

import (
	"bufio"
	"encoding/json"
	"github.com/easmith/p2p-messanger/proto"
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

	lastPeers := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastPeers = append(lastPeers, scanner.Text())
	}

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

	// TODO: добавить в пиры

	// TODO: ping-pong
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
