package proto

type PeerName struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

type PeerList struct {
	Cmd   string     `json:"cmd"`
	Peers []PeerName `json:"peers"`
}

type WsCmd struct {
	Cmd     string `json:"cmd"`
	To      string `json:"to"`
	Content string `json:"content"`
}
