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

		decodedMessage := &proto.WsCmd{}
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
		if decodedMessage.Cmd == "MESS" {

			hexPubKey, err := hex.DecodeString(decodedMessage.To)
			if err != nil {
				log.Printf("decode error: %s", err)
				continue
			}
			peer, found := p.Peers.Get(string(hexPubKey))
			if found {
				writeToWs(c, mt, message)
				p.SendMessage(*peer.Conn, decodedMessage.Content)
			}
		}
	}
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
