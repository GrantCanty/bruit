package kraken_data

import (
	"log"
)

var (
	ApiKey     string
	PrivateKey string
)

func LoadKeys(env map[string]string) {
	ApiKey = env["API_KEY"]
	PrivateKey = env["PRIVATE_KEY"]

	log.Println("ApiKey: ", ApiKey, "PrivateKey: ", PrivateKey)
}
