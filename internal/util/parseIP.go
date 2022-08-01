package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type MinecraftAddress struct {
	IP       string
	Port     uint16
	Combined string
}

func ParseIP(raw string) (addr MinecraftAddress, err error) {
	fmt.Println("parsing IP....")
	fmt.Println(net.ParseIP(raw))

	if strings.Contains(raw, ":") {
		split := strings.Split(raw, ":")
		port, err := strconv.Atoi(split[1])
		if err != nil {
			return addr, err
		}

		addr.IP = split[0]
		addr.Port = uint16(port)
	} else {
		addr.IP = raw
		addr.Port = 25565
	}

	addr.Combined = fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	return
}
