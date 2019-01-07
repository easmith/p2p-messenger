package listener

import (
	"bufio"
	"fmt"
	"github.com/easmith/p2p-messenger/proto"
	"github.com/easmith/p2p-messenger/types"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

func StartListener(port int, ch chan string, proto *proto.Proto) {
	service := fmt.Sprintf(":%v", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Printf("ResolveTCPAddr: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("ListenTCP: %s", err.Error())
		os.Exit(1)
	}

	log.Println("Start listen " + service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			ch <- "conn Accept error: " + err.Error()
			continue
		}
		// TODO: общение через канал
		ch <- "new connection: " + conn.RemoteAddr().String()
		go onConnection(conn, proto)
	}

	ch <- "done"
}

func onConnection(conn net.Conn, p *proto.Proto) {
	defer func() {
		//proto.Peers.(conn)
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
	log.Printf("Stop read peak: %v", string(buf))

	if types.ItIsHttp(buf) {
		log.Println("Try request")
		handleHttp(readWriter, conn, p)
	} else {
		log.Println("Try proto")
		proto.HandleProto(readWriter, conn, p)
	}
}

func handleHttp(rw *bufio.ReadWriter, conn net.Conn, p *proto.Proto) {
	request, err := http.ReadRequest(rw.Reader)

	if err != nil {
		log.Printf("Read request ERROR: %s", err)
		return
	}

	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	s := conn.RemoteAddr().String()[0:3] + "REMOVE_IT"
	// TODO: сравнение среза со строкой
	if strings.EqualFold(s, "127") || strings.EqualFold(s, "[::") {
		response.Body = ioutil.NopCloser(strings.NewReader("php-messenger 1.0"))
	} else {

		if path.Clean(request.URL.Path) == "/ws" {
			handleWs(NewWriter(conn), request, p)
			return
		} else {
			processRequest(request, &response)
			//fileServer := http.FileServer(http.Dir("./front/build/"))
			//fileServer.ServeHTTP(NewWriter(conn), request)
		}
	}

	err = response.Write(rw)
	if err != nil {
		log.Printf("Write response ERROR: %s", err)
		return
	}

	err = rw.Writer.Flush()
	if err != nil {
		log.Printf("Flush response ERROR: %s", err)
		return
	}
}
