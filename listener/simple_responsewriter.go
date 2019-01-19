package listener

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

//MyWriter Простейшая реализация интерфейса ResponseWriter
type MyWriter struct {
	conn net.Conn
}

//Write ResponseWriter implementation
func (w MyWriter) Write(b []byte) (int, error) {
	return w.conn.Write(b)
}

//Header ResponseWriter implementation
func (w MyWriter) Header() http.Header {
	return http.Header{}
}

//WriteHeader ResponseWriter implementation
func (w MyWriter) WriteHeader(statusCode int) {
	_, err := w.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK")))
	if err != nil {
		log.Printf("WriteHeaderError: %v\n", err)
	}
}

//Hijack захват сокета. Используется при апгрейде соединения до WebSocket
func (w MyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	reader := bufio.NewReader(w.conn)
	writer := bufio.NewWriter(w.conn)

	readWriter := bufio.NewReadWriter(reader, writer)
	return w.conn, readWriter, nil
}

//NewMyWriter create new ResponseWriter implementation
func NewMyWriter(conn net.Conn) http.ResponseWriter {
	return &MyWriter{conn}
}
