package proto

import "encoding/json"

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

type PeerName struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

func (v PeerName) ToJson() []byte {
	return toJson(v)
}

type HandShake struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
	ExKey  string `json:"exKey"`
}

func (v HandShake) ToJson() []byte {
	return toJson(v)
}

type WsCmd struct {
	Cmd string `json:"cmd"`
}

type WsMyName struct {
	WsCmd
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

func (v WsMyName) ToJson() []byte {
	return toJson(v)
}

type WsPeerList struct {
	WsCmd
	Peers []PeerName `json:"peers"`
}

func (v WsPeerList) ToJson() []byte {
	return toJson(v)
}

type WsMessage struct {
	WsCmd
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

func (v WsMessage) ToJson() []byte {
	return toJson(v)
}
