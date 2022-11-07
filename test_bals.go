package main

import (
	"bruit/bruit"
	"bruit/bruit/clients/kraken"
)

func main() {
	settings := bruit.Settings{}
	settings.Init()

	k := kraken.KrakenClient{}
	k.InitClient(&settings)
}
