package config

import "os"

var (
	CacheMode bool
)

func LoadConfigVars() {
	CacheMode = os.Getenv("CACHE") == "true"
}
