/* ThreadedEchoServer
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"local/p2pmessager/types"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
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

	service := fmt.Sprintf(":%v", *port)

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

	go SocketClient("", 1)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn)
	}

}

func SocketClient(ip string, port int) {
	//addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", ":35035")

	defer conn.Close()

	if err != nil {
		log.Fatalln("asdasd" + err.Error())
	}

	conn.Write([]byte("HAND" + name))

	buff := make([]byte, 128)
	n, err := conn.Read(buff)
	log.Printf("Receive(%v): %s, %v\n", n, buff[:n], err)

}

func onConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readWriter := bufio.NewReadWriter(reader, writer)

	buf, err := readWriter.Peek(4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
		return
	}

	if types.ItIsMessanger(buf) {
		log.Println("This is message!")
		handleMessage(readWriter)
	} else {
		log.Println("Try request")
		handleRequest(readWriter)
	}
}

func handleMessage(rw *bufio.ReadWriter) {
	var buf = make([]byte, 128)
	for {
		n, err := rw.Read(buf)
		if err != nil {
			return
		}
		log.Printf("MAIN Recieved: [%v] %s (%s)", n, buf, err)

		nn, err := rw.Write([]byte("NAMEOK"))
		log.Printf("MAIN Sent: [%v] (%s)", nn, err)
		rw.Writer.Flush()

		time.Sleep(100000)
	}
}

func handleRequest(rw *bufio.ReadWriter) {
	request, err := http.ReadRequest(rw.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err)
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

type MyWriter struct {
	*bufio.Writer
}

func (w MyWriter) Header() http.Header {
	return http.Header{}
}

func (w MyWriter) WriteHeader(statusCode int) {
	w.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK")))
}

func NewWriter(w *bufio.Writer) http.ResponseWriter {
	return &MyWriter{w}
}
