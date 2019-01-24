package proto

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"log"
)

var cmdLen = 4
var idLen = 16
var fromLen = 32
var toLen = 32
var signLen = 64
var headerLen = cmdLen + idLen + fromLen + toLen + signLen + 2

//Envelope Конверт для сообщений между пирами
type Envelope struct {
	Cmd     []byte
	Id      []byte
	From    []byte
	To      []byte
	Sign    []byte
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
func NewEnvelope(cmd string, contentBytes []byte) (envelope *Envelope) {
	contentLength := len(contentBytes)
	if contentLength >= 65535 {
		contentBytes = contentBytes[:65535]
	}

	envelope = &Envelope{
		Cmd:     []byte(cmd)[:cmdLen],
		Id:      getRandomSeed(idLen)[:idLen],
		From:    make([]byte, fromLen),
		To:      make([]byte, toLen),
		Sign:    make([]byte, signLen),
		Length:  uint16(contentLength),
		Content: contentBytes[0:contentLength],
	}
	return
}

//NewSignedEnvelope create new envelop with signature
func NewSignedEnvelope(cmd string, from []byte, to []byte, sign []byte, contentBytes []byte) (envelope *Envelope) {
	envelope = NewEnvelope(cmd, contentBytes)
	envelope.From = from
	envelope.To = to
	envelope.Sign = sign
	return
}

//Serialize Сериализация конверта и содержимого в массив байт
func (m Envelope) Serialize() []byte {
	result := make([]byte, 0, headerLen+len(m.Content))

	// TODO: неудобная конкатенация
	result = append(result, m.Cmd[0:cmdLen]...)
	result = append(result, m.Id[0:idLen]...)
	result = append(result, m.From[0:fromLen]...)
	result = append(result, m.To[0:toLen]...)
	result = append(result, m.Sign[0:signLen]...)

	contentLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLengthBytes, m.Length)

	result = append(result, contentLengthBytes...)
	result = append(result, m.Content...)

	return result
}

//UnSerialize Десериализация массива байт в конверт с содержимым
func UnSerialize(b []byte) (envelope *Envelope) {
	contentLength := binary.BigEndian.Uint16(b[headerLen-2 : headerLen])
	if contentLength > 65535 {
		return nil
	}

	envelope = &Envelope{
		Cmd:    b[0:cmdLen],
		Id:     b[cmdLen : cmdLen+idLen],
		From:   b[cmdLen+idLen : cmdLen+idLen+fromLen],
		To:     b[cmdLen+idLen+fromLen : cmdLen+idLen+fromLen+toLen],
		Sign:   b[cmdLen+idLen+fromLen+toLen : cmdLen+idLen+fromLen+toLen+signLen],
		Length: contentLength,
	}

	if len(b) == (headerLen + int(contentLength)) {
		envelope.Content = b[headerLen:]
	} else {
		envelope.Content = make([]byte, contentLength)
	}

	return
}

//ReadEnvelope Формирование конверта из байтов ридера сокета
func ReadEnvelope(reader *bufio.Reader) (*Envelope, error) {
	header := make([]byte, headerLen)

	// read envelope header
	_, err := reader.Read(header)
	if err != nil {
		return nil, err
	}

	envelope := UnSerialize(header)

	_, err = reader.Read(envelope.Content)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

//Send send envelop to peer
func (m Envelope) Send(peer *Peer) {
	log.Printf("Send %s to peer %s ", m.Cmd, peer.Name)
	_, err := (*peer.Conn).Write(m.Serialize())
	if err != nil {
		log.Printf("ERROR on write message: %v", err)
	}
}
