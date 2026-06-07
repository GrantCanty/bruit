package main

import (
	"bruit/bruit/client"
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken"
	"log"
)

func main() {
	var c clients.BruitCryptoClient
	c = &kraken.KrakenClient{}

	var bruit client.BruitClient
	bruit = &client.Client{}
	if err := bruit.Init(c); err != nil {
		log.Println("Failed to init bruit: ", err)
		return
	}
	go bruit.Run()

	bruit.Wait()
}
