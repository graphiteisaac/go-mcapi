package minecraft

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
)

func PingServer(addr *Address) (string, error) {
	_, srvs, err := net.LookupSRV("minecraft", "tcp", addr.Host)
	if err == nil {
		sort.Slice(srvs, func(i, j int) bool {
			return srvs[i].Weight > srvs[j].Weight
		})

		srv := srvs[0]
		addr.Port = srv.Port
		addr.Host = srv.Target
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr.Host, addr.Port), 10*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// send packet to server
	var dataBuf bytes.Buffer
	var finBuf bytes.Buffer

	dataBuf.Write([]byte("\x00"))                 // start packet
	dataBuf.Write([]byte("\x6D"))                 // protocol (-1 default)
	dataBuf.Write([]uint8{uint8(len(addr.Host))}) // length of host
	dataBuf.Write([]byte(addr.Host))              // host
	a := make([]byte, 2)                          // port
	binary.BigEndian.PutUint16(a, addr.Port)      // port
	dataBuf.Write(a)                              // port
	dataBuf.Write([]byte("\x01"))                 // end

	pacLen := []byte{uint8(dataBuf.Len())}
	finBuf.Write(append(pacLen, dataBuf.Bytes()...))

	conn.Write(finBuf.Bytes())
	conn.Write([]byte("\x01\x00"))

	// read packet response
	read := bufio.NewReader(conn)
	binary.ReadUvarint(read)

	packetType, _ := read.ReadByte()
	if bytes.Compare([]byte{packetType}, []byte("\x00")) != 0 {
		return "", errors.New("response packet byte mismatch")
	}

	//Get data length via Varint
	length, err := binary.ReadUvarint(read)
	if err != nil {
		return "", err
	}

	if length < 10 {
		return "", errors.New("response too small")
	} else if length > 700000 {
		return "", errors.New("response too large")
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
