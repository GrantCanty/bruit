package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"encoding/json"
	"log"
	"strconv"
)

// PUBLIC SOCKET METHODS

func (client *KrakenClient) SubscribeToTrades(g settings.Settings, pairs []string) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	if err := PubSocketGuard(client.WebSocket); err != nil {
		panic(err)
	}

	client.WebSocket.SubscribeToTrades(pairs)
}

/****
	*Add func to check if already subscribed to OHLC Stream
	*Add func to get past OHLC data from rest API. Add to the candle map list
*****/
func (client *KrakenClient) SubscribeToOHLC(g settings.Settings, pairs []string, interval int) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	var found bool = false
	for _, i := range kraken_data.OHLCVals {
		if i == interval {
			found = true
			break
		}
	}
	if found == true {
		if err := PubSocketGuard(client.WebSocket); err != nil { // guard clause checker
			panic(err)
		}

		// add func here that makes request to rest OHLC to get past OHLC data. data should then be added to the OHLC map

		client.WebSocket.SubscribeToOHLC(pairs, interval)
	}
}

func (client *KrakenClient) PubDecoder(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

	if err := PubSocketGuard(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocketPointer())
	client.WebSocket.GetPubSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PubJsonDecoder(message, g.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocketPointer())

	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
	return
}

// ORDER BOOK SOCKET METHODS

// Subscribe to the order book.
func (client *KrakenClient) SubscribeToOrderBook(g settings.Settings, pairs []string, depth int) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	if err := BookSocketGuard(client.WebSocket); err != nil {
		panic(err)
	}

	sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.BookSubscription{
			Depth: depth,
			Name:  "book",
		},
		Pair: pairs,
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
	log.Println(string(sub))
	client.WebSocket.GetBookSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) BookDecoder(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

	if err := BookSocketGuard(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetBookSocketPointer())
	client.WebSocket.GetBookSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.BookJsonDecoder(message, g.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocketPointer())

	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
	return
}

// PRIVATE SOCKET METHODS

func (client *KrakenClient) SubscribeToOpenOrders(g settings.Settings, token string) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	PrivSocketGuard(client.WebSocket)

	/*sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name:  "openOrders",
			Token: token,
		},
	})

	if err != nil {
		panic(err)
	}

	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)*/
	client.WebSocket.SubscribeToOpenOrders(token)
}

func (client *KrakenClient) CancelAll(g settings.Settings, token string) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	PrivSocketGuard(client.WebSocket)

	sub, _ := json.Marshal(&types.Subscribe{
		Event: "cancelAll",
		Token: token,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)
}

// find a way to ad tradeID
func (client *KrakenClient) CancelOrder(g settings.Settings, token string, tradeIDs []string) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	PrivSocketGuard(client.WebSocket)

	sub, _ := json.Marshal(&types.CancelOrder{
		Event: "cancelOrder",
		Token: token,
		Txid:  tradeIDs,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) AddOrder(g settings.Settings, token string, otype string, ttype string, pair string, vol string, price string, testing bool) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	PrivSocketGuard(client.WebSocket)

	test := strconv.FormatBool(testing)
	sub, _ := json.Marshal(&types.Order{
		WsToken:   token,
		Event:     "addOrder",
		OrderType: otype,
		TradeType: ttype,
		Pair:      pair,
		Volume:    vol,
		Price:     price,
		Validate:  test,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) PrivDecoder(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

	PrivSocketGuard(client.WebSocket)

	ws_client.ReceiveLocker(client.WebSocket.GetPrivSocketPointer())
	client.WebSocket.GetPrivSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PrivJsonDecoder(message, g.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPrivSocketPointer())

	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
	return
}
