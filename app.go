package main

import (
	"flag"
	"github.com/easmith/p2p-messenger/discover"
	"github.com/easmith/p2p-messenger/listener"
	"github.com/easmith/p2p-messenger/proto"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	// Настройки логирования
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	name := flag.String("name", "You", "name")
	port := flag.Int("port", 35035, "port as port")
	peersFile := flag.String("peers", "peers.txt", "peers file")

	flag.Parse()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("Exit by signal: %s", sig)
		os.Exit(1)
	}()

	p := proto.NewProto(*name, *port)

	var wg sync.WaitGroup

	wg.Add(2)
	go discover.StartDiscover(p, *peersFile)
	go listener.StartListener(p, *port)
	wg.Wait()

	//e := webview.Open("Peer To Peer Messenger", fmt.Sprintf("http://localhost:%v", *port), 800, 600, false)
	//if e != nil {
	//	log.Printf(e.Error())
	//}
}
