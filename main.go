package main

import (
	"bruit/bruit/client"
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken"
)

func main() {
	var c clients.BruitCryptoClient
	c = &kraken.KrakenClient{}

	var bruit client.BruitClient
	bruit = &client.Client{}
	bruit.Init(c)
	go bruit.Run()

	bruit.Wait()
}
