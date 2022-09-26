package decoders

import (
	"bruit_new/bruit/clients/kraken/types"
	"encoding/json"
	"log"
)

func OpenOrdersResponseDecoder(byteResponse []byte, testing bool) (*types.OpenOrdersResponse, error) {
	if testing == true {
		log.Println("in openOrdersResponse func")
	}

	var resp types.OpenOrdersResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("openOrdersResponse error: ", err)
		}
		return nil, err
	}
	return &resp, err
}

func CancelOrderResponseDecoder(byteResponse []byte, testing bool) (*types.CancelOrderResponse, error) {
	if testing == true {
		log.Println("in open orders response func")
	}

	var resp types.CancelOrderResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("open orders response error: ", err)
		}
		return nil, err
	}
	return &resp, err
}
