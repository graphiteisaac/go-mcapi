package minecraft

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

type pingResponsePlayers struct {
	Online int `json:"online"`
	Max    int `json:"max"`
}

type pingResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PingResponse struct {
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

func MarshalJson(raw string) (*PingResponse, error) {
	ping := PingResponse{}

	if err := json.Unmarshal([]byte(raw), &ping); err != nil {
		return nil, err
	}

	return &ping, nil
}
