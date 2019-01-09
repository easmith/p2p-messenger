package proto

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"log"
	"net"
)

var idLen = 16
var cmdLen = 4

type Message struct {
	Cmd     []byte
	MsgId   []byte
	Length  uint16
	Content []byte
}

func getRandomSeed(l int) []byte {
	seed := make([]byte, l)
	_, err := rand.Read(seed)
	if err != nil {
		log.Printf("rand.Read Error: %v", err)
	}
	return seed
}

func NewMessage(cmd string, contentBytes []byte) Message {
	contentLength := len(contentBytes)
	log.Printf("content[%v]: %s", contentLength, contentBytes)
	if contentLength >= (2 << 16) {
		contentLength = (2 << 16) - 1
		contentBytes = contentBytes[:contentLength]
	}

	return Message{
		Cmd:     []byte(cmd)[:cmdLen],
		MsgId:   getRandomSeed(idLen)[:idLen],
		Length:  uint16(contentLength),
		Content: contentBytes[0:contentLength],
	}
}

func (m Message) Serialize() []byte {
	result := make([]byte, 0, cmdLen+idLen+len(m.Content))

	// TODO: неудобная конкатенация
	result = append(result, m.Cmd[0:cmdLen]...)
	result = append(result, m.MsgId[0:idLen]...)

	contentLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLengthBytes, m.Length)

	result = append(result, contentLengthBytes...)
	result = append(result, m.Content...)

	return result
}

//UnSerialize
func UnSerialize(b []byte) Message {
	contentLength := binary.BigEndian.Uint16(b[idLen+cmdLen : idLen+cmdLen+2])

	return Message{
		Cmd:     b[0:cmdLen],
		MsgId:   b[cmdLen : cmdLen+idLen],
		Length:  contentLength,
		Content: b[idLen+cmdLen+2 : idLen+cmdLen+2+int(contentLength)],
	}
}

func ReadMessage(reader *bufio.Reader) (*Message, error) {
	msgId := make([]byte, idLen)
	cmd := make([]byte, cmdLen)
	contentLength := make([]byte, 2)
	_, err := reader.Read(cmd)
	// TODO: Много лишнего кода при обработке ошибок
	if err != nil {
		return nil, err
	}
	_, err = reader.Read(msgId)
	if err != nil {
		return nil, err
	}
	_, err = reader.Read(contentLength)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(contentLength)
	content := make([]byte, length)
	_, err = reader.Read(content)
	if err != nil {
		return nil, err
	}
	return &Message{
		MsgId:   msgId,
		Cmd:     cmd,
		Length:  length,
		Content: content,
	}, nil
}

func (m Message) WriteToConn(conn net.Conn) {
	log.Printf("Proto write: %s", m.Cmd)
	_, err := conn.Write(m.Serialize())
	if err != nil {
		log.Printf("ERROR on write message: %v", err)
	}
}
