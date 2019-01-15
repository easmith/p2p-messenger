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

// Переключение на WebSocket и обмен сообщений с фронтом через него
func handleWs(w http.ResponseWriter, r *http.Request, p *proto.Proto) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	go waitMessageForWs(p, c)

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

		switch decodedMessage.Cmd {
		case "PEERS":
			{
				peerList := p.Peers.PeerList()

				peerListJson, err := json.Marshal(peerList)

				if err != nil {
					panic(err)
				}

				writeToWs(c, mt, peerListJson)
			}
		case "MESS":
			{
				hexPubKey, err := hex.DecodeString(decodedMessage.To)
				if err != nil {
					log.Printf("decode error: %s", err)
					continue
				}
				peer, found := p.Peers.Get(string(hexPubKey))
				if found {
					writeToWs(c, mt, message)
					p.SendMessage(peer, decodedMessage.Content)
				}

			}
		}

	}
}

func waitMessageForWs(p *proto.Proto, c *websocket.Conn) {
	for {
		envelope := <-p.Broker
		log.Printf("New message: %s", envelope.Cmd)
		if string(envelope.Cmd) == "MESS" {

			wsCmd := proto.WsMessage{
				WsCmd: proto.WsCmd{
					Cmd: "MESS",
				},
				To:      "ME",
				Content: string(envelope.Content),
			}

			wsCmdBytes, err := json.Marshal(wsCmd)

			if err != nil {
				panic(err)
			}

			writeToWs(c, 1, wsCmdBytes)
		}
	}
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
