package main

import (
	"flag"
	"github.com/easmith/p2p-messenger/discover"
	"github.com/easmith/p2p-messenger/listener"
	"github.com/easmith/p2p-messenger/proto"
	"log"
	"os"
)

func main() {

	// Настройки логирования
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	name := flag.String("name", "ONE", "name")
	port := flag.Int("port", 35035, "port as port")

	flag.Parse()

	// Устновака порта в случае неправильного ввода
	if *port <= 0 || *port > 65535 {
		*port = 35035
	}

	waitClose := make(chan int)

	proto := proto.NewProto(*name)

	go discover.StartDiscover(proto)

	go listener.StartListener(*port, proto)

	for {
		log.Printf("Close: %v", <-waitClose)
	}
}
