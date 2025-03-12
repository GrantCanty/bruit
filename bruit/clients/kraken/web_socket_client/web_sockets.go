package web_socket

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	decoders "bruit/bruit/clients/kraken/decoder_funcs"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"encoding/json"
	"log"
)

type MessageTypeIdentifier struct {
    Type string `json:"type"`
}

type WebSocketClient struct {
	pubSocket  ws_client.Socket
	bookSocket ws_client.Socket
	privSocket ws_client.Socket

	bookChan     chan interface{} // this chan contains the final book data
	bookJSONChan chan interface{} // this chan contains the most recently decoded data from the book subscription
	privChan     chan interface{}
}

func (client *WebSocketClient) PubJsonDecoder(response string, logger settings.LoggingSettings, OHLCch chan types.OHLCResponse, Tradech chan types.TradeResponse, OHLCsubch chan types.OHLCSuccessResponse) {
	byteResponse := []byte(response)

	if ohlcResp, err := decoders.OhlcResponseDecoder(byteResponse, logger.GetLoggingConsole()); err == nil {
        OHLCch <- *ohlcResp
        return
    }
	if tradeResp, err := decoders.TradeResponseDecoder(byteResponse, logger.GetLoggingConsole()); err == nil {
        Tradech <- *tradeResp
        return
    }
	if _, err := decoders.HbResponseDecoder(byteResponse, logger.GetLoggingConsole()); err == nil {
		return
	}
	if _, err := decoders.ServerConnectionStatusResponseDecoder(byteResponse, logger.GetLoggingConsole()); err == nil {
		return
	}
	if ohlcSubResp, err := decoders.OhlcSubscriptionResponseDecoder(byteResponse, logger.GetLoggingConsole()); err == nil {
		OHLCsubch <- *ohlcSubResp
		return
	}

	return
}

func (client *WebSocketClient) BookJsonDecoder(response string, logger settings.LoggingSettings) {
	byteResponse := []byte(response)

	var msgType MessageTypeIdentifier
    if err := json.Unmarshal(byteResponse, &msgType); err != nil {
        log.Println("Error identifying message type:", err)
        return
    }

	switch msgType.Type {
    case "update":
        if resp, err := decoders.UpdateBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
            log.Println(resp)
        } else {
            log.Println("Error decoding update:", err)
        }
    case "snapshot":
        if resp, err := decoders.SnapshotBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
            log.Println(resp)
        } else {
            log.Println("Error decoding snapshot:", err)
        }
    default:
        log.Println("Unknown message type:", msgType.Type)
    }

	/*if resp, err := decoders.UpdateBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
		log.Println(resp)
		return
	}
	if resp, err := decoders.SnapshotBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
		log.Println(resp)
		return
	}*/
	
	//if resp, err = 
	/*var resp interface{}
	byteResponse := []byte(response)

	resp, err := decoders.InitialBookResponseDecoder(byteResponse, now, logger.GetLoggingConsole())
	if err != nil {
		resp, err = decoders.IncrementalAskAndBidDecoder(byteResponse, logger.GetLoggingConsole())
		if err != nil {
			resp, err = decoders.IncrementalAskOrBidDecoder(byteResponse, logger.GetLoggingConsole())
			if err != nil {
				resp, err = decoders.HbResponseDecoder(byteResponse, logger.GetLoggingConsole())
				if err != nil {
					resp, err = decoders.ServerConnectionStatusResponseDecoder(byteResponse, logger.GetLoggingConsole())
					if err != nil {
						resp, err = decoders.BookSubscriptionResponseDecoder(byteResponse, logger.GetLoggingConsole())
						if err != nil {
							log.Println(string("\033[31m"), "Received response of unknown data type: ", response)
						}
					}
				}
			}
		}
	}*/
	/*log.Println("resp from boonJsonDecoder: ", resp)
	byteResponse := []byte(response)

	err := json.Unmarshal(byteResponse, &resp)
	if err != nil {
		log.Fatal(err)
	}

	for _, element := range resp[1].(map[string]interface{})["as"].([]interface{}) {
		priceStr := element.([]interface{})[0].(string)
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			log.Fatal(err)
		}
		volStr := element.([]interface{})[1].(string)
		vol, err := decimal.NewFromString(volStr)
		if err != nil {
			log.Fatal(err)
		}
		ob.Asks = append(ob.Asks, types.Level{Price: price, Volume: vol})
	}

	for _, element := range resp[1].(map[string]interface{})["bs"].([]interface{}) {
		priceStr := element.([]interface{})[0].(string)
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			log.Fatal(err)
		}
		volStr := element.([]interface{})[1].(string)
		vol, err := decimal.NewFromString(volStr)
		if err != nil {
			log.Fatal(err)
		}
		ob.Bids = append(ob.Bids, types.Level{Price: price, Volume: vol})
	}*/
	//client.bookJSONChan <- resp
	return
}

func (client *WebSocketClient) PrivJsonDecoder(response string, logger settings.LoggingSettings) {
	var resp interface{}
	byteResponse := []byte(response)

	resp, err := decoders.OpenOrdersResponseDecoder(byteResponse, logger.GetLoggingConsole())
	if err != nil {
		resp, err = decoders.HbResponseDecoder(byteResponse, logger.GetLoggingConsole())
		if err != nil {
			resp, err = decoders.CancelOrderResponseDecoder(byteResponse, logger.GetLoggingConsole())
			if err != nil {
				resp, err = decoders.ServerConnectionStatusResponseDecoder(byteResponse, logger.GetLoggingConsole())
			}
		}
	}

	//if testing == true {
	//	log.Println(reflect.TypeOf(resp), resp)
	//}
	client.privChan <- resp
	return //resp

	/*var resp interface{}
	byteResponse := []byte(response)

	resp, err := decoders.OhlcResponseDecoder(byteResponse, l.GetLoggingConsole()) // these funcs need to accept LoggingSettings struct so they can take both DBlogging and ConsoleLogging
	if err != nil {
		resp, err = decoders.HbResponseDecoder(byteResponse, l.GetLoggingConsole())
		if err != nil {
			resp, err = decoders.ServerConnectionStatusResponseDecoder(byteResponse, l.GetLoggingConsole())
			if err != nil {
				resp, err = decoders.OhlcSubscriptionResponseDecoder(byteResponse, l.GetLoggingConsole())
				if err != nil {
					log.Println(string("\033[31m"), "Received response of unknown data type: ", response)
				}
			}
		}
	}

	client.pubChan <- resp
	return*/
}

func (ws *WebSocketClient) InitChannels() {
	//ws.pubChan = make(chan interface{})
	ws.bookChan = make(chan interface{})
	ws.bookJSONChan = make(chan interface{})
	ws.privChan = make(chan interface{})
}

func (ws *WebSocketClient) SubscribeToTrades(pairs []string) {
	sub, _ := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name: "trade",
		},
		Pair: pairs,
	})
	ws.pubSocket.SendBinary(sub)
}

func (ws *WebSocketClient) SubscribeToOHLC(pairs []string, interval int) {
	sub, _ := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.OHLCSubscription{
			Interval: interval,
			Name:     "ohlc",
		},
		Pair: pairs,
	})
	ws.pubSocket.SendBinary(sub)
}

func (ws *WebSocketClient) SubscribeToOpenOrders(token string) {
	sub, _ := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name:  "openOrders",
			Token: token,
		},
	})
	ws.privSocket.SendBinary(sub)
}

func (ws *WebSocketClient) SubscribeToOrderBook(pairs []string, depth int) {
	log.Println(pairs)
	sub, err := json.Marshal(&types.SubscribeV2{
		Method: "subscribe",
		Params: types.ParamsV2{
			Depth: depth,
			Channel:  "book",
			Symbol: pairs,
		},
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
	ws.bookSocket.SendBinary(sub)
}

func (client *WebSocketClient) InitSockets() { // used to initialized public and private sockets
	ws_client.PublicInit(&client.pubSocket, kraken_data.GetPubWSUrl())
	ws_client.BookInit(&client.bookSocket, kraken_data.GetV2WsURL())
	ws_client.PrivateInit(&client.privSocket, kraken_data.GetPrivWSUrl())
}

func (ws WebSocketClient) GetPubSocket() ws_client.Socket {
	return ws.pubSocket
}

func (ws *WebSocketClient) GetPubSocketPointer() *ws_client.Socket {
	return &ws.pubSocket
}

func (ws WebSocketClient) GetBookSocket() ws_client.Socket {
	return ws.bookSocket
}

func (ws *WebSocketClient) GetBookSocketPointer() *ws_client.Socket {
	return &ws.bookSocket
}

func (ws WebSocketClient) GetPrivSocket() ws_client.Socket {
	return ws.privSocket
}

func (ws *WebSocketClient) GetPrivSocketPointer() *ws_client.Socket {
	return &ws.privSocket
}

func (ws *WebSocketClient) GetBookChan() chan interface{} {
	return ws.bookChan
}

func (ws *WebSocketClient) GetBookJSONChan() chan interface{} {
	return ws.bookJSONChan
}

func (ws *WebSocketClient) GetPrivChan() chan interface{} {
	return ws.privChan
}
