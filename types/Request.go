package types

import (
	"fmt"
	"strings"
)

type Request struct {
	Verb    string
	Headers map[string]string
	Body    string
}


func Parse(ba []byte) (request Request) {
	request.Headers = make(map[string]string)
	split := strings.Split(string(ba), "\r\n")

	request.Verb = split[0]

	for _, v := range split[1:] {
		kv := strings.Split(v, ":")
		if len(kv) > 1 {
			request.Headers[kv[0]] = kv[1]
		}
	}

	return
}

func (r Request) String() string {
	result := r.Verb + "\r\n"

	for k, v := range r.Headers {
		result += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	result += "\r\n" + r.Body
	return result
}
