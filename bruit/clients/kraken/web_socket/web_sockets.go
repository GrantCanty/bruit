package web_socket

import (
	"bruit/bruit"
	kraken_data "bruit/bruit/clients/kraken/client_data"
	"bruit/bruit/clients/kraken/decoders"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/ws_client"
	"encoding/json"
	"log"
)

type WebSocketClient struct {
	pubSocket  ws_client.Socket
	bookSocket ws_client.Socket
	privSocket ws_client.Socket

	pubChan  chan interface{}
	bookChan chan interface{}
	privChan chan interface{}

	//subscriptions map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData
}

/*func (ws WebSocketClient) GetSubscriptions() map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData {
	return ws.subscriptions
}

func (ws WebSocketClient) GetInterval(metaData shared_types.SubscriptionMetaData) int {
	return ws.subscriptions[metaData].GetData().(types.KrakenOHLCSubscriptionData).Interval
}

func (ws WebSocketClient) GetChannelID(metaData shared_types.SubscriptionMetaData) int {
	return metaData.GetChannelID()
}

func (ws *WebSocketClient) HandleOHLCSuccessResponse(resp *types.OHLCSuccessResponse) {
	if _, found := ws.subscriptions[resp.GetMetaData()]; found {
		ws.subscriptions[resp.GetMetaData()] = types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status}
	} else {
		if ws.subscriptions == nil {
			ws.subscriptions = make(map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData)
		}
		ws.subscriptions[resp.GetMetaData()] = types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status}
	}
	log.Println("subscription list: ", ws.subscriptions)
}*/

func (client *WebSocketClient) PubJsonDecoder(response string, logger bruit.LoggingSettings) {
	var resp interface{}
	byteResponse := []byte(response)

	resp, err := decoders.OhlcResponseDecoder(byteResponse, logger.GetLoggingConsole()) // these funcs need to accept LoggingSettings struct so they can take both DBlogging and ConsoleLogging
	if err != nil {
		resp, err = decoders.TradeResponseDecoder(byteResponse, logger.GetLoggingConsole())
		if err != nil {
			resp, err = decoders.HbResponseDecoder(byteResponse, logger.GetLoggingConsole())
			if err != nil {
				resp, err = decoders.ServerConnectionStatusResponseDecoder(byteResponse, logger.GetLoggingConsole())
				if err != nil {
					resp, err = decoders.OhlcSubscriptionResponseDecoder(byteResponse, logger.GetLoggingConsole())
					if err != nil {
						log.Println(string("\033[31m"), "Received response of unknown data type: ", response)
					}
				}
			}
		}
	}

	client.pubChan <- resp
	return
}

func (client *WebSocketClient) BookJsonDecoder(response string, logger bruit.LoggingSettings) {
	var resp BookResp
	//byteResponse := []byte(response)

	client.bookChan <- resp
}

func (client *WebSocketClient) PrivJsonDecoder(response string, logger bruit.LoggingSettings) {
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
	ws.pubChan = make(chan interface{})
	ws.bookChan = make(chan interface{})
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

func (client *WebSocketClient) InitConnections() { // used to initialized public and private sockets
	ws_client.PublicInit(&client.pubSocket, kraken_data.GetPubWSUrl())
	ws_client.BookInit(&client.bookSocket, kraken_data.GetPubWSUrl())
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

func (ws *WebSocketClient) GetPubChan() chan interface{} {
	return ws.pubChan
}

func (ws *WebSocketClient) GetBookChan() chan interface{} {
	return ws.bookChan
}

func (ws *WebSocketClient) GetPrivChan() chan interface{} {
	return ws.privChan
}
