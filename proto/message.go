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

//Envelope Конверт для сообщений между пирами
type Envelope struct {
	Cmd     []byte
	Id      []byte
	Length  uint16
	Content []byte
}

func (m Envelope) String() string {
	return string(m.Cmd) + "-" + string(m.Id) + "-" + string(m.Length)
}

func getRandomSeed(l int) []byte {
	seed := make([]byte, l)
	_, err := rand.Read(seed)
	if err != nil {
		log.Printf("rand.Read Error: %v", err)
	}
	return seed
}

//NewEnvelope Создание нового конверта
func NewEnvelope(cmd string, contentBytes []byte) Envelope {
	contentLength := len(contentBytes)
	log.Printf("content[%v]: %s", contentLength, contentBytes)
	if contentLength >= (2 << 16) {
		contentLength = (2 << 16) - 1
		contentBytes = contentBytes[:contentLength]
	}

	return Envelope{
		Cmd:     []byte(cmd)[:cmdLen],
		Id:      getRandomSeed(idLen)[:idLen],
		Length:  uint16(contentLength),
		Content: contentBytes[0:contentLength],
	}
}

//Serialize Сериализация конверта и содержимого в массив байт
func (m Envelope) Serialize() []byte {
	result := make([]byte, 0, cmdLen+idLen+len(m.Content))

	// TODO: неудобная конкатенация
	result = append(result, m.Cmd[0:cmdLen]...)
	result = append(result, m.Id[0:idLen]...)

	contentLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLengthBytes, m.Length)

	result = append(result, contentLengthBytes...)
	result = append(result, m.Content...)

	return result
}

//UnSerialize Десериализация массива байт в конверт с содержимым
func UnSerialize(b []byte) Envelope {
	contentLength := binary.BigEndian.Uint16(b[idLen+cmdLen : idLen+cmdLen+2])

	return Envelope{
		Cmd:     b[0:cmdLen],
		Id:      b[cmdLen : cmdLen+idLen],
		Length:  contentLength,
		Content: b[idLen+cmdLen+2 : idLen+cmdLen+2+int(contentLength)],
	}
}

//ReadEnvelope Формирование конверта из байтов ридера сокета
func ReadEnvelope(reader *bufio.Reader) (*Envelope, error) {
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
	return &Envelope{
		Id:      msgId,
		Cmd:     cmd,
		Length:  length,
		Content: content,
	}, nil
}

func (m Envelope) WriteToConn(conn net.Conn) {
	log.Printf("Proto write: %s", m.Cmd)
	_, err := conn.Write(m.Serialize())
	if err != nil {
		log.Printf("ERROR on write message: %v", err)
	}
}
