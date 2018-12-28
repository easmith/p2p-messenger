package proto

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
	"log"
	"net"
	"os"
)

type Addr string

type Peer struct {
	PubKey    ed25519.PublicKey
	Addr      Addr
	Conn      *net.Conn
	Name      string
	FirstSeen string
	LastSeen  string
	Peers     *Peers
}

func (p Peer) String() string {
	return p.Name + "=" + hex.EncodeToString(p.PubKey)
}

type Proto struct {
	Name    string
	Peers   *Peers
	PubKey  ed25519.PublicKey
	privKey ed25519.PrivateKey
}

func getSeed() []byte {
	seed := getRandomSeed(32)

	fName := "seed.dat"
	file, err := os.Open(fName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(fName)
			if err != nil {
				panic(err)
			}
		}
	}

	_, err = file.Read(seed)
	if err != nil {
		panic(err)
	}
	return seed
}

func NewProto(name string) Proto {
	//privateKey := ed25519.NewKeyFromSeed(getSeed())
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return Proto{
		Name:    name,
		Peers:   NewPeers(),
		PubKey:  publicKey,
		privKey: privateKey,
	}
}

type PeerName struct {
	Name   string
	PubKey string
}

func (p Proto) SendName(conn net.Conn) {
	peerName, err := json.Marshal(PeerName{
		Name:   p.Name,
		PubKey: hex.EncodeToString(p.PubKey),
	})

	if err != nil {
		panic(err)
	}
	message := NewMessage("NAME", peerName)
	message.WriteToConn(conn)
}

func (p Proto) RequestPeers(conn net.Conn) {
	message := NewMessage("LIST", []byte("TODO"))
	message.WriteToConn(conn)
}

func (p Proto) SendPeers(conn net.Conn) {
	message := NewMessage("PEER", []byte("TODO"))
	message.WriteToConn(conn)
}

func (p Proto) SendMessage(conn net.Conn, msg string) {
	message := NewMessage("MESS", []byte(msg))
	message.WriteToConn(conn)
}

func ConnListener(conn net.Conn, p *Proto) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	HandleProto(readWriter, conn, p)
}

func HandleProto(rw *bufio.ReadWriter, conn net.Conn, p *Proto) {
	for {
		message, err := ReadMessage(rw.Reader)
		if err != nil {
			log.Printf("Error on read Message: %v", err)
			return
		}

		log.Printf("new Message: %s %s", message.Cmd, message.Content)

		if string(message.Cmd) == "NAME" {
			peerName := PeerName{}
			err := json.Unmarshal(message.Content, &peerName)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			log.Printf("recieve name: %v", peerName)
			p.SendName(conn)
		}
	}
}
