package listener

import (
	"bufio"
	"fmt"
	"github.com/easmith/p2p-messanger/proto"
	"github.com/easmith/p2p-messanger/types"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

func onConnection(conn net.Conn, proto *proto.Proto) {
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
		handleHttp(readWriter, conn)
	} else {
		log.Println("Try proto")
		handleProto(readWriter, conn, proto)
	}
}

func handleProto(rw *bufio.ReadWriter, conn net.Conn, p *proto.Proto) {
	for {
		message, err := proto.ReadMessage(rw.Reader)
		if err != nil {
			log.Printf("Error on read Message: %v", err)
			continue
		}

		log.Printf("new Message: %v %v %v", message.MsgId, message.Cmd, message.Length)
	}
}

func handleHttp(rw *bufio.ReadWriter, conn net.Conn) {
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
			handleWs(NewWriter(conn), request)
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

func handleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		log.Printf("ws read: [%v] %s", mt, message)
		err = c.WriteMessage(mt, append([]byte("server recv: "), message...))
		if err != nil {
			log.Printf("ws write error: %s", err)
			break
		}
	}
}
