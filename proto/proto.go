package proto

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/ed25519"
	"log"
	"net"
	"os"
	"reflect"
)

//Proto Ядро протокола
type Proto struct {
	Name    string
	Peers   *Peers
	PubKey  ed25519.PublicKey
	privKey ed25519.PrivateKey
	Broker  chan Envelope
}

func (p Proto) String() string {
	return "proto: " + hex.EncodeToString(p.PubKey) + ": " + p.Name
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

//NewProto - создание экземпляра ядра протокола
func NewProto(name string) *Proto {
	//privateKey := ed25519.NewKeyFromSeed(getSeed())
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return &Proto{
		Name:    name,
		Peers:   NewPeers(),
		PubKey:  publicKey,
		privKey: privateKey,
		Broker:  make(chan Envelope),
	}
}

//SendName Отправка своего имени в сокет
func (p Proto) SendName(conn net.Conn) {
	peerName, err := json.Marshal(PeerName{
		Name:   p.Name,
		PubKey: hex.EncodeToString(p.PubKey),
	})

	if err != nil {
		panic(err)
	}
	message := NewEnvelope("NAME", peerName)
	message.WriteToConn(conn)
}

//RequestPeers Запрос списка пиров
func (p Proto) RequestPeers(conn net.Conn) {
	message := NewEnvelope("LIST", []byte("TODO"))
	message.WriteToConn(conn)
}

//SendPeers Отправка списка пиров
func (p Proto) SendPeers(conn net.Conn) {
	message := NewEnvelope("PEER", []byte("TODO"))
	message.WriteToConn(conn)
}

//SendMessage Отправка сообщения
func (p Proto) SendMessage(conn net.Conn, msg string) {
	message := NewEnvelope("MESS", []byte(msg))
	message.WriteToConn(conn)
}

//RegisterPeer Регистрация пира в списках пиров
func (p Proto) RegisterPeer(peer *Peer) *Peer {
	// TODO: сравнение через equal
	if reflect.DeepEqual(peer.PubKey, p.PubKey) {
		return nil
	}

	p.Peers.Put(peer)

	log.Printf("Register new peer: %s", peer.Name)

	return peer
}

//UnregisterPeer Удаление пира из списка
func (p Proto) UnregisterPeer(peer *Peer) {
	p.Peers.Remove(peer)
	log.Printf("UnRegister peer: %s", peer.Name)
}

//PeerListener Старт прослушивания соединения с пиром
func (p Proto) PeerListener(peer *Peer) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(*peer.Conn), bufio.NewWriter(*peer.Conn))
	p.HandleProto(readWriter, *peer.Conn)
}

//HandleProto Обработка входящих сообщений
func (p Proto) HandleProto(rw *bufio.ReadWriter, conn net.Conn) {
	var peer *Peer
	for {
		message, err := ReadEnvelope(rw.Reader)
		if err != nil {
			log.Printf("Error on read Envelope: %v", err)
			break
		}

		switch string(message.Cmd) {
		case "NAME":
			{
				newPeer := CreatePeer(message, conn)
				if newPeer != nil {
					if peer != nil {
						p.UnregisterPeer(peer)
					}
					p.RegisterPeer(newPeer)
					peer = newPeer
				}
				p.SendName(conn)
			}
		case "MESS":
			{
				log.Printf("NEW MESSAGE %s", message.Content)
			}
		default:
			log.Printf("PROTO MESSAGE %v %v %v", message.Cmd, message.Id, message.Content)
		}

	}

	if peer != nil {
		p.UnregisterPeer(peer)
	}

}
