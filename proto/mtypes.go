package proto

import "encoding/json"

//Serializable interface to detect that can to serialised to json
type Serializable interface {
	ToJson() []byte
}

func toJson(v interface{}) []byte {
	json, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return json
}

//PeerName Peer name and public key
type PeerName struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

//ToJson convert to JSON bytes
func (v PeerName) ToJson() []byte {
	return toJson(v)
}

//HandShake type for handshake on connection
type HandShake struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
	ExKey  string `json:"exKey"`
}

//ToJson convert to JSON bytes
func (v HandShake) ToJson() []byte {
	return toJson(v)
}

//WsCmd WebSocket command
type WsCmd struct {
	Cmd string `json:"cmd"`
}

//WsMyName WebSocket command: PeerName
type WsMyName struct {
	WsCmd
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

//ToJson convert to JSON bytes
func (v WsMyName) ToJson() []byte {
	return toJson(v)
}

//WsPeerList WebSocket command: list of peers
type WsPeerList struct {
	WsCmd
	Peers []PeerName `json:"peers"`
}

//ToJson convert to JSON bytes
func (v WsPeerList) ToJson() []byte {
	return toJson(v)
}

//WsMessage WebSocket command: new Message
type WsMessage struct {
	WsCmd
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

//ToJson convert to JSON bytes
func (v WsMessage) ToJson() []byte {
	return toJson(v)
}
