package minecraft

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// TODO come up with better way of handling descriptions, maybe build HTML?
//type pingResponseDescription struct {
//	Text  string      `json:"text"`
//	Extra interface{} `json:"extra,omitempty"`
//}

type pingResponsePlayers struct {
	Online int `json:"online"`
	Max    int `json:"max"`
}

type pingResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PingResponse struct {
	Hostname    string              `json:"hostname"`
	IPv4        string              `json:"ipv4"`
	Port        uint16              `json:"port"`
	Cached      bool                `json:"cached"`
	Description interface{}         `json:"description"`
	Players     pingResponsePlayers `json:"players"`
	Version     pingResponseVersion `json:"version"`
	Icon        string              `json:"favicon"`
}

func GetIcon(raw string) ([]byte, error) {
	var icon struct {
		Icon string `json:"favicon"`
	}

	err := json.Unmarshal([]byte(raw), &icon)
	if err != nil {
		return nil, err
	}

	image := icon.Icon[strings.IndexByte(icon.Icon, ',')+1:]

	return base64.StdEncoding.DecodeString(image)
}

func CreateResponseObj(raw string, host string, port uint16, cached bool) (ping PingResponse, err error) {
	err = json.Unmarshal([]byte(raw), &ping)

	ping.Hostname = host
	ping.Port = port
	ping.Cached = cached
	ping.IPv4 = host
	ping.Hostname = host

	return
}
