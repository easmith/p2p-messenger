package types

var itHttp = map[string]bool{
	"GET ": true,
	"HEAD": true,
	"POST": true,
}

var itMessanger = map[string]bool{
	"NAME": true,
	"LIST": true,
	"SEND": true,
}

func ItIsHttp(ba []byte) bool {
	return itHttp[string(ba)]
}

func ItIsMessanger(ba []byte) bool {
	return itMessanger[string(ba)]
}
