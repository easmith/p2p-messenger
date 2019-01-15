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

//StartListener Старт прослушивания порта и обработки входящих соединений
func StartListener(proto *proto.Proto, port int) {
	// Устновака порта в случае неправильного ввода
	if port <= 0 || port > 65535 {
		port = 35035
	}

	service := fmt.Sprintf("localhost:%v", port)

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

	fmt.Printf("\n\tService start on %s\n\n", tcpAddr.String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn, proto)
	}

}

// Обработка входящего соединения
func onConnection(conn net.Conn, p *proto.Proto) {
	defer func() {
		//proto.Peers.(conn)
		conn.Close()
	}()

	log.Printf("New connection from: %v", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readWriter := bufio.NewReadWriter(reader, writer)

	buf, err := readWriter.Peek(4)
	if err != nil {
		log.Printf("Read peak ERROR: %s", err)
		return
	}

	if types.ItIsHttp(buf) {
		handleHttp(readWriter, conn, p)
	} else {
		p.HandleProto(readWriter, conn)
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
		response.Body = ioutil.NopCloser(strings.NewReader("Peer To Peer Messenger. see https://github.com/easmith/p2p-messenger"))
	} else {

		if path.Clean(request.URL.Path) == "/ws" {
			handleWs(NewMyWriter(conn), request, p)
			return
		} else {
			processRequest(request, &response)
			//fileServer := http.FileServer(http.Dir("./front/build/"))
			//fileServer.ServeHTTP(NewMyWriter(conn), request)
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
