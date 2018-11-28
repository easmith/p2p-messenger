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
	for {
		n, err := rw.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Disconected by EOF")
			} else {
				log.Printf("handlePeer ERROR: %s", err)
			}
			return
		}

		log.Printf("MAIN Recieved: [%v] %s", n, bytes.Trim(buf, "\r\n\x00"))

		switch string(buf[0:4]) {
		case "NAME":
			{
				peer := peers.Add(conn, types.Id(string(bytes.Trim(buf[4:16], "\r\n\x00"))))
				log.Printf("new peer: %s", peer)
				conn.Write([]byte("REGI\n"))
			}
		case "LIST":
			{
				for _, peer := range peers.ById.Peers {
					conn.Write([]byte(fmt.Sprintf("PEER\t%v\t%s\n", peer.Addr, peer.Id)))
				}

			}
		default:
			conn.Write([]byte("UNKNOWN_CMD\n"))
		}

		if string(buf[0:4]) == "LIST" {

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
