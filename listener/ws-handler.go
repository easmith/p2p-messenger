package listener

import (
	"encoding/hex"
	"encoding/json"
	"github.com/easmith/p2p-messenger/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWs(w http.ResponseWriter, r *http.Request, p *proto.Proto) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		log.Printf("ws read: [%v] %s", mt, message)

		decodedMessage := &proto.WsMessage{}
		err = json.Unmarshal(message, decodedMessage)

		if err != nil {
			log.Printf("error on unmarshal message: %v", err)
			continue
		}

		if decodedMessage.Cmd == "PEERS" {

			peerList := p.Peers.PeerList()

			peerListJson, err := json.Marshal(peerList)

			if err != nil {
				panic(err)
			}

			writeToWs(c, mt, peerListJson)
		}
		if string(message[0:4]) == "MESS" {

			hexPubKey, err := hex.DecodeString(string(message[4:68]))
			if err != nil {
				log.Printf("LISTENER: decode error: %s", err)
				continue
			}
			peer, found := p.Peers.Get(string(hexPubKey))
			if found {
				p.SendMessage(*peer.Conn, string(message[68:]))
			}
		}
	}
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, append([]byte("server recv: "), message...))
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
