package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mc-api/internal/config"
	"mc-api/internal/db"
	"mc-api/internal/util"
	"net"
	"strings"
	"time"
)

func PingServerTCP(c context.Context, addr util.MinecraftAddress) (ping util.PingResponse, err error) {
	conn, err := net.DialTimeout("tcp", addr.Combined, 10*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	// send packet to server
	util.SendPacket(&conn, addr.IP, addr.Port)

	// read packet response
	res, err := util.ReadPacketResponse(&conn)
	if err != nil {
		return ping, errors.New("tcp error: cant read packet response")
	}

	// unmarshal response into ping obj
	ping, err = util.CreateResponseObj(res, addr, false)
	fmt.Println(conn.RemoteAddr())
	ping.IPv4 = strings.Split(conn.RemoteAddr().String(), ":")[0]

	if err != nil {
		fmt.Println(err)
		return ping, errors.New("malformed response: cannot unmarshal response")
	}

	b, err := json.Marshal(ping)

	if config.CacheMode {
		err = db.Redis.Set(c, addr.Combined, string(b), time.Minute*3).Err()
		if err != nil {
			return ping, errors.New("redis error: cannot set")
		}
	}

	return
}
