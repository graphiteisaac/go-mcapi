package minecraft

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

func PingServer(addr Address) (ping PingResponse, err error) {
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
		return
	}
	defer conn.Close()

	// send packet to server
	SendPacket(&conn, addr.IP, addr.Port)

	// read packet response
	res, err := ReadPacketResponse(&conn)

	if err != nil {
		return ping, errors.New("tcp error: cant read packet response")
	}

	// unmarshal response into ping obj
	ping, err = CreateResponseObj(res, addr.IP, addr.Port, false)
	ping.IPv4 = strings.Split(conn.RemoteAddr().String(), ":")[0]

	if err != nil {
		fmt.Println(err)
		return ping, errors.New("malformed response: cannot unmarshal response")
	}

	return ping, nil
}
