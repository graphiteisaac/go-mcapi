package util

import (
	"encoding/base64"
	"strings"
)

func Base64StringToImage(input string) ([]byte, error) {
	image := input[strings.IndexByte(input, ',')+1:]

	return base64.StdEncoding.DecodeString(image)
}
