package util

import (
	"encoding/json"
)

type pingResponseDescription struct {
	Text  string      `json:"text"`
	Extra interface{} `json:"extra,omitempty"`
}

type pingResponsePlayers struct {
	Online int `json:"online"`
	Max    int `json:"max"`
}

type pingResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PingResponse struct {
	Hostname    string                  `json:"hostname"`
	IPv4        string                  `json:"ipv4"`
	Port        uint16                  `json:"port"`
	Cached      bool                    `json:"cached"`
	Description pingResponseDescription `json:"description"`
	Players     pingResponsePlayers     `json:"players"`
	Version     pingResponseVersion     `json:"version"`
	Icon        string                  `json:"favicon"`
}

func CreateResponseObj(raw string, address MinecraftAddress, cached bool) (ping PingResponse, err error) {
	ping.Hostname = address.IP
	ping.Port = address.Port
	ping.Cached = cached

	err = json.Unmarshal([]byte(raw), &ping)

	return
}
