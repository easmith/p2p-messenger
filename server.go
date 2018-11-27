/* ThreadedEchoServer
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"local/p2pmessager/types"
	"log"
	"net"
	"os"
)

var name string

func main() {

	name = *flag.String("name", "defaultName", "name")
	port := flag.Int("port", 35035, "port as port")

	flag.Parse()

	print(name)

	if *port <= 0 || *port > 65535 {
		*port = 35035
	}

	service := fmt.Sprintf(":%v", *port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service )
	if err != nil {
		fmt.Fprintf(os.Stderr, "ResolveTCPAddr: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListenTCP: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Start listen %s\n", service)

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
		log.Fatalln(err)
	}

	conn.Write([]byte("HAND" + name))

	//buff := make([]byte, 1024)
	//_, _ := conn.Read(buff)
	//log.Printf("Receive: %s", buff[:n])

}


func onConnection(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
listener:
	for {
		_, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		head := buf[:4]
		if types.ItIsHttp(head) {
			//readRequest, e := http.ReadRequest(buf)

			request := types.Parse(buf[0:])
			fmt.Printf("%s\n", request.Verb)
			fmt.Printf("%s\n", request)
			writePage(conn)
			break listener
		}

		if types.ItIsMessanger(head) {
			switch {
			case bytes.Equal([]byte("HAND"), head) : {
				handshake(buf[4:16], conn)
			}
			case bytes.Equal([]byte("CLOS"), head) : {
				break listener
			}
			}
		}
	}
}

func writePage(conn net.Conn) {
	response := types.New("hello!")

	fmt.Printf("Response\n%s\n", response)

	_, err := fmt.Fprintf(conn, "%s", response)
	if err != nil {
		return
	}
}

func handshake(handName []byte, conn net.Conn) {
	fmt.Printf("Handshake from: %s", handName)
	_, err := fmt.Fprintf(conn, "HAND%s", name)
	if err != nil {
		return
	}
}