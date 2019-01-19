/*
Package to discovering new peers on network and to announcing yourself
*/
package discover

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/easmith/p2p-messenger/proto"
	"log"
	"net"
	"os"
	"time"
)

var peers = make(map[string]string)

//StartDiscover Начинает подключения к пирам из списка peers.txt и отправляет им свое имя
func StartDiscover(p *proto.Proto, peersFile string) {

	go startMeow("224.0.0.1:35035", p)
	go listenMeow("224.0.0.1:35035", p, connectToPeer)

	file, err := os.Open(peersFile)
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
		go connectToPeer(p, peerAddress)
	}
}

// подключаемся к пиру по адресу
func connectToPeer(p *proto.Proto, peerAddress string) {
	if _, exist := peers[peerAddress]; exist {
		log.Printf("peer %s already exist", peerAddress)
		return
	}
	peers[peerAddress] = peerAddress
	log.Printf("try to connect peer: %s", peerAddress)

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

	p.RegisterPeer(peer)

	p.ListenPeer(peer)

	p.UnregisterPeer(peer)

	delete(peers, peerAddress)
	// TODO: ping-pong
}

// Отправка UPD multicast пакетов с информацией о себе
func startMeow(address string, p *proto.Proto) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Printf(err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf(err.Error())
	}

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("meow:%v:%v", hex.EncodeToString(p.PubKey), p.Port)))
		if err != nil {
			log.Printf(err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

// Прослушка UPD
func listenMeow(address string, p *proto.Proto, handler func(p *proto.Proto, peerAddress string)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetReadBuffer(1024)
	if err != nil {
		log.Fatal(err)
	}

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, 1024)
		_, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		trim := bytes.Trim(buffer, "\x00")

		peerPubKeyStr := string(trim[5 : 5+64])
		peerPubKey, err := hex.DecodeString(peerPubKeyStr)

		// Если с этим пиром уже есть связь, то пропускаем его
		_, found := p.Peers.Get(string(peerPubKey))
		if found || bytes.Equal(p.PubKey, peerPubKey) {
			continue
		}

		peerAddress := src.IP.String() + string(trim[5+64:])

		log.Printf(" peer: %v %v", peerPubKeyStr, found)
		log.Printf("Found peer: %s", peerAddress)

		handler(p, peerAddress)
	}
}

// Отправляем пиру свое имя и ожидаем от него его имя
func handShake(p *proto.Proto, conn net.Conn) *proto.Peer {
	log.Printf("DISCOVERY: try handshake with %s", conn.RemoteAddr())
	peer := proto.NewPeer(conn)

	p.SendName(peer)

	envelope, err := proto.ReadEnvelope(bufio.NewReader(conn))
	if err != nil {
		log.Printf("Error on read Envelope: %s", err)
		return nil
	}

	err = peer.UpdatePeer(envelope)
	if err != nil {
		log.Printf("HandShake error: %s", err)
	}

	return peer
}

func iface() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	var addrs = make(map[string]net.Interface)
	for _, i := range ifaces {
		multicastAddrs, err := i.MulticastAddrs()
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		for _, a := range multicastAddrs {
			addrs[a.String()] = i
			log.Printf("%v\n", a.String())
		}
	}

	log.Printf("=====")
	for k, _ := range addrs {
		log.Printf("%v", k)
	}
}
