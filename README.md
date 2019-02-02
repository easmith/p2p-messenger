# Peer to Peer messaging

[![GoDoc](https://godoc.org/github.com/easmith/p2p-messenger?status.svg)](https://godoc.org/github.com/easmith/p2p-messenger)
[![Go Report Card](https://goreportcard.com/badge/github.com/easmith/p2p-messenger)](https://goreportcard.com/report/github.com/easmith/p2p-messenger)
[![LICENSE](https://img.shields.io/github/license/easmith/p2p-messenger.svg)](https://github.com/easmith/p2p-messenger/blob/master/LICENSE)



Build front with simple UI:
```bash
cd front
npm update
npm run build
```
    
Start messaging:
```bash 
cd ..
go run app.go -name Snowden
```
Start params

    -name string
        you name (default "USER@HOSTNAME")
    -peers string
        Path to file with peer addresses on each line (default "peers.txt")
    -port int
        port that have to listen (default 35035)
    -webview
        Start WebView ui (default true)
        
Extended info
- https://easmith.github.io/post/golang-p2p-messenger/
- https://habr.com/ru/post/437686/