package services

import (
	"context"
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
	ping.IPv4 = strings.Split(conn.RemoteAddr().String(), ":")[0]

	if err != nil {
		fmt.Println(err)
		return ping, errors.New("malformed response: cannot unmarshal response")
	}

	if config.CacheMode {
		// TODO modify redis expiry to find something clean
		err = db.Redis.Set(c, addr.Combined, res, time.Minute*3).Err()
		if err != nil {
			return ping, errors.New("redis error: cannot set")
		}
	}

	return
}
