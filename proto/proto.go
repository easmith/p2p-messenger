package proto

import (
	"golang.org/x/crypto/ed25519"
	"net"
	"os"
)

type Addr string
type PubKey string

type Peer struct {
	Id        PubKey
	Addr      Addr
	Conn      *net.Conn
	FirstSeen string
	LastSeen  string
	Peers     map[PubKey]*Peer
}

type Proto struct {
	Name    string
	Peers   map[PubKey]*Peer
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
		Peers:   map[PubKey]*Peer{},
		PubKey:  publicKey,
		privKey: privateKey,
	}
}

func (p Proto) SendName(conn net.Conn) {
	message := NewMessage("NAME", "myName")
	message.WriteToConn(conn)
}

func (p Proto) RequestPeers(conn net.Conn) {
	message := NewMessage("LIST", "")
	message.WriteToConn(conn)
}

func (p Proto) SendPeers(conn net.Conn) {
	message := NewMessage("PEER", "")
	message.WriteToConn(conn)
}
