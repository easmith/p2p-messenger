package listener

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"local/p2pmessager/types"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

func StartListener(port int, ch chan string, peers *types.Peers) {
	service := fmt.Sprintf(":%v", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ResolveTCPAddr: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListenTCP: %s", err.Error())
		os.Exit(1)
	}

	log.Println("Start listen " + service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			ch <- "conn Accept error: " + err.Error()
			continue
		}
		ch <- "new connection: " + conn.RemoteAddr().String()
		go onConnection(conn, peers)
	}

	ch <- "done"
}

func onConnection(conn net.Conn, peers *types.Peers) {
	defer func() {
		peers.Remove(conn)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readWriter := bufio.NewReadWriter(reader, writer)

	buf, err := readWriter.Peek(4)
	log.Printf("Start read peak")
	if err != nil {
		log.Printf("Read peak ERROR: %s", err)
		return
	}
	log.Printf("Stop read peak")

	if types.ItIsMessanger(buf) {
		log.Println("This is message!")
		handlePeer(readWriter, conn, peers)
	} else {
		log.Println("Try request")
		handleRequest(readWriter)
	}
}

func handlePeer(rw *bufio.ReadWriter, conn net.Conn, peers *types.Peers) {
	var buf = make([]byte, 1024)
	var peer *types.Peer
	for {
		var cmd = make([]byte, 4)
		n, err := rw.Read(cmd)
		if err != nil {
			if err == io.EOF {
				log.Printf("Disconected by EOF")
			} else {
				log.Printf("handlePeer ERROR: %s", err)
			}
			return
		}

		log.Printf("MAIN Recieved: [%v] %s", n, bytes.Trim(cmd, "\r\n\x00"))

		switch string(cmd) {
		case "NAME":
			{
				line, _, _ := rw.ReadLine()
				id := types.Id(bytes.Trim(line, " "))
				_, found := peers.ById.Get(id)
				if found {
					conn.Write([]byte("ERR Name already in use\n"))
					continue
				}
				peer = peers.Add(&conn, id)
				log.Printf("new peer: %s", peer)
				conn.Write([]byte("OK\n"))
			}
		case "LIST":
			{
				rw.Read(buf)
				for _, p := range peers.ById.Peers {
					conn.Write([]byte(fmt.Sprintf("PEER\t%v\t%s\n", p.Addr, p.Id)))
				}

			}
		case "SEND":
			{
				if peer == nil {
					//conn.Write([]byte(fmt.Sprintf("ERR you not registered (use name)\n")))
					continue
				}

				to, _ := rw.ReadString(0x20)

				p, found := peers.ById.Get(types.Id(to[0 : len(to)-1]))
				if !found {
					n, e := conn.Write([]byte(fmt.Sprintf("ERR not found %v\n", types.Id(to))))
					continue
				}

				msg, err := rw.ReadString('\n')
				if err != nil {
					log.Printf("ReadMessageError: %s", err)
					continue
				}

				_, err = (*p.Conn).Write([]byte(fmt.Sprintf("MESS %v: %v", peer.Id, msg)))
				if err != nil {
					log.Printf("WriteMessageError: %s", err)
					(*peer.Conn).Write([]byte("ERR " + err.Error()))
					continue
				}

				(*peer.Conn).Write([]byte("OK\n"))
			}
		default:
			conn.Write([]byte("UNKNOWN_CMD\n"))
		}
	}
}

func handleRequest(rw *bufio.ReadWriter) {
	request, err := http.ReadRequest(rw.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read request ERROR: %s", err)
		return
	}

	log.Printf("Request: %v\t%v\n", request.URL.Path, path.Clean(request.URL.Path))

	//fs := http.FileServer(http.Dir("/home/eku/go/src/local/p2pmessager/static/"))

	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	file, e := os.Open("./static" + path.Clean(request.URL.Path))
	if e != nil {
		log.Printf("error: %s", e)
		response.StatusCode = 404
		response.Body = ioutil.NopCloser(strings.NewReader("Not found"))
		response.Write(rw)
		rw.Writer.Flush()
		return
	}

	response.Body = file
	response.Write(rw)

	rw.Writer.Flush()
}

//type MyWriter struct {
//	*bufio.Writer
//}
//
//func (w MyWriter) Header() http.Header {
//	return http.Header{}
//}
//
//func (w MyWriter) WriteHeader(statusCode int) {
//	w.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK")))
//}
//
//func NewWriter(w *bufio.Writer) http.ResponseWriter {
//	return &MyWriter{w}
//}
//
