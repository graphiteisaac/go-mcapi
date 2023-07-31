package minecraft

import (
	"strconv"
	"strings"
)

type Address struct {
	Host string
	Port uint16
}

func ParseIP(raw string) (*Address, error) {
	addr := &Address{
		Host: raw,
		Port: 25565,
	}

	if strings.Contains(raw, ":") {
		split := strings.Split(raw, ":")
		port, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}

		addr.Host = split[0]
		addr.Port = uint16(port)
	}

	return addr, nil
}
