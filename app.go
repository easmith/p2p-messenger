package main

import (
	"flag"
	"github.com/easmith/p2p-messanger/discover"
	"github.com/easmith/p2p-messanger/listener"
	"github.com/easmith/p2p-messanger/proto"
	"log"
	"os"
)

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	name := flag.String("name", "ONE", "name")
	port := flag.Int("port", 35035, "port as port")

	flag.Parse()

	if *port <= 0 || *port > 65535 {
		*port = 35035
	}

	proto := proto.NewProto(*name)

	listenerChan := make(chan string)

	go discover.StartDiscover(&proto)

	go listener.StartListener(*port, listenerChan, &proto)

	for {
		log.Printf("Message from listener channel: %s", <-listenerChan)
	}
}
