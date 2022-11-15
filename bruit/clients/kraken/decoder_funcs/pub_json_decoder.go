package decoders

import (
	"bruit/bruit/clients/kraken/types"
	"bytes"
	"encoding/json"
	"log"
)

func OhlcResponseDecoder(byteResponse []byte, testing bool) (*types.OHLCResponse, error) {
	if testing == true {
		log.Println("in ohlcResponse func")
	}

	var resp types.OHLCResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("ohlcResponse error: ", err)
		}
		return nil, err
	}
	return &resp, err
}

func TradeResponseDecoder(byteResponse []byte, testing bool) (*types.TradeResponse, error) {
	if testing == true {
		log.Println("in tradeResponse func")
	}

	var resp types.TradeResponse
	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		if testing == true {
			log.Println("tradeResponse error: ", err)
		}
		return nil, err
	}
	return &resp, err
}

func HbResponseDecoder(byteResponse []byte, testing bool) (*types.HeartBeat, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in hb response func")
	}
	var heart types.HeartBeat
	err := decoder.Decode(&heart)
	if err != nil {
		if testing == true {
			log.Println(err)
		}
		return nil, err
	}
	return &heart, err
}

func ServerConnectionStatusResponseDecoder(byteResponse []byte, testing bool) (*types.ServerConnectionStatusResponse, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in connection status func")
	}
	var conn types.ServerConnectionStatusResponse
	err := decoder.Decode(&conn)
	if err != nil {
		if testing == true {
			log.Println("connection status error: ", err)
		}
		return nil, err
	}
	return &conn, err
}

func OhlcSubscriptionResponseDecoder(byteResponse []byte, testing bool) (*types.OHLCSuccessResponse, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in ohlc subscription response func")
	}
	var ohlc types.OHLCSuccessResponse
	err := decoder.Decode(&ohlc)
	if err != nil {
		if testing == true {
			log.Println("ohlc subscription response error", err)
		}
		return nil, err
	}
	return &ohlc, err
}
