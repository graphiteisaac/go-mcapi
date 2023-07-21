package minecraft

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
)

func PingServer(addr Address) (string, error) {
	_, srvs, err := net.LookupSRV("minecraft", "tcp", addr.IP)
	if err == nil {
		sort.Slice(srvs, func(i, j int) bool {
			return srvs[i].Weight > srvs[j].Weight
		})

		srv := srvs[0]
		addr.Combined = fmt.Sprintf("%v:%d", srv.Target, srv.Port)
		addr.Port = srv.Port
		addr.IP = srv.Target
	}

	conn, err := net.DialTimeout("tcp", addr.Combined, 10*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// send packet to server
	SendPacket(&conn, addr.IP, addr.Port)

	// read packet response
	res, err := ReadPacketResponse(&conn)
	if err != nil {
		return res, errors.New("tcp error: cant read packet response")
	}

	if err != nil {
		fmt.Println(err)
		return res, errors.New("malformed response: cannot unmarshal response")
	}

	return res, nil
}
