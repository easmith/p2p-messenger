package types

var itHttp = map[string]bool{
	"GET ": true,
	"HEAD": true,
	"POST": true,
}

func ItIsHttp(ba []byte) bool {
	return itHttp[string(ba)]
}
