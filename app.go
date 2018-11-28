/* ThreadedEchoServer
 */
package main

import (
	"flag"
	"local/p2pmessager/listener"
	"local/p2pmessager/types"
	"log"
	"os"
)

var name string

func main() {

	log.SetOutput(os.Stdout)

	name = *flag.String("name", "ONE", "name")
	port := flag.Int("port", 35035, "port as port")

	flag.Parse()

	if *port <= 0 || *port > 65535 {
		*port = 35035
	}

	peers := types.NewPeers()

	listenerChan := make(chan string)

	//go discover.Start("", 1)

	go listener.StartListener(*port, listenerChan, peers)

	for {
		log.Printf("Message from listener channel: %s", <-listenerChan)
	}
}
