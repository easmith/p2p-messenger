/*
Go training project
*/
package main

import (
	"flag"
	"fmt"
	"github.com/easmith/p2p-messenger/discover"
	"github.com/easmith/p2p-messenger/listener"
	"github.com/easmith/p2p-messenger/proto"
	"github.com/zserge/webview"
	"log"
	"os"
	"os/signal"
	"os/user"
	"sync"
	"syscall"
)

//InitParams initialising params
type InitParams struct {
	Name         string
	Port         int
	PeersFile    string
	StartWebView bool
}

var initParams InitParams

func init() {
	currentUser, _ := user.Current()
	hostName, _ := os.Hostname()

	initParams = InitParams{
		Name:         *flag.String("name", currentUser.Username+"@"+hostName, "you name"),
		Port:         *flag.Int("port", 35035, "port that have to listen"),
		PeersFile:    *flag.String("peers", "peers.txt", "Path to file with peer addresses on each line"),
		StartWebView: *flag.Bool("webview", true, "Start WebView ui"),
	}

	flag.Parse()

	flag.PrintDefaults()

	// Настройки логирования
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

}

func main() {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("Exit by signal: %s", sig)
		os.Exit(1)
	}()

	p := proto.NewProto(initParams.Name, initParams.Port)

	if initParams.StartWebView {
		startWithWebView(p)
	} else {
		startWithoutWebView(p)
	}
}

//startWithWebView Запуск приложения с WebView
func startWithWebView(p *proto.Proto) {
	go discover.StartDiscover(p, initParams.PeersFile)
	go listener.StartListener(p, initParams.Port)

	if initParams.StartWebView {
		e := webview.Open("Peer To Peer Messenger", fmt.Sprintf("http://localhost:%v", initParams.Port), 800, 600, false)
		if e != nil {
			log.Printf(e.Error())
		}
	}
}

//startWithoutWebView Запуск приложения без запуска WebView
func startWithoutWebView(p *proto.Proto) {
	var wg sync.WaitGroup
	wg.Add(2)
	go discover.StartDiscover(p, initParams.PeersFile)
	go listener.StartListener(p, initParams.Port)
	wg.Wait()
}
