package minecraft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

func ReadPacketResponse(conn *net.Conn) (string, error) {
	read := bufio.NewReader(*conn)
	binary.ReadUvarint(read)

	packetType, _ := read.ReadByte()
	if bytes.Compare([]byte{packetType}, []byte("\x00")) != 0 {
		return "", errors.New("error response packet type")
	}

	//Get data length via Varint
	length, err := binary.ReadUvarint(read)
	if err != nil {
		return "", err
	}

	if length < 10 {
		return "", errors.New("error to small response")
	} else if length > 700000 {
		return "", errors.New("error to big response")
	}

	// Receive json buffer
	bytesReceived := uint64(0)
	recBytes := make([]byte, length)
	for bytesReceived < length {
		n, _ := read.Read(recBytes[bytesReceived:length])
		bytesReceived = bytesReceived + uint64(n)
	}

	return string(recBytes), nil
}

func SendPacket(conn *net.Conn, host string, port uint16) {
	var dataBuf bytes.Buffer
	var finBuf bytes.Buffer

	dataBuf.Write([]byte("\x00"))            // start packet
	dataBuf.Write([]byte("\x6D"))            // protocol (-1 default)
	dataBuf.Write([]uint8{uint8(len(host))}) // length of host
	dataBuf.Write([]byte(host))              // host
	a := make([]byte, 2)                     // port
	binary.BigEndian.PutUint16(a, port)      // port
	dataBuf.Write(a)                         // port
	dataBuf.Write([]byte("\x01"))            // end

	pacLen := []byte{uint8(dataBuf.Len())}
	finBuf.Write(append(pacLen, dataBuf.Bytes()...))

	(*conn).Write(finBuf.Bytes())
	(*conn).Write([]byte("\x01\x00"))
}