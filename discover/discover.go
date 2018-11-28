package discover

import (
	"log"
	"net"
)

func StartDiscover(name string) {

	var conn net.Conn
	var err error

	if conn, err = net.Dial("tcp", ":35035"); err != nil {
		log.Printf("Dial Error: " + err.Error())
	}

	defer conn.Close()

	conn.Write([]byte("HAND" + "NaMe"))

	buff := make([]byte, 128)
	n, err := conn.Read(buff)
	log.Printf("Receive(%v): %s, %v\n", n, buff[:n], err)

}
