package util

import (
	"fmt"
	"net"
)

type MinecraftAddress struct {
	IP       string
	Port     uint16
	Combined string
}

func ParseIP(raw string) (addr MinecraftAddress, err error) {
	fmt.Println(net.ParseIP(raw))
	addr.IP = raw
	addr.Port = 25565
	addr.Combined = fmt.Sprintf("%s:%d", addr.IP, addr.Port)

	return
}
