package proto

type PeerName struct {
	Name   string `json:"name"`
	PubKey string `json:"id"`
}

type WsCmd struct {
	Cmd string `json:"cmd"`
}

type PeerList struct {
	WsCmd
	Peers []PeerName `json:"peers"`
}

type WsMessage struct {
	WsCmd
	To      string `json:"to"`
	Content string `json:"content"`
}
