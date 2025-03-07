package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

func remove(slice []string, pos int) []string {
	return append(slice[:pos], slice[pos+1:]...)
}

// PUBLIC SOCKET METHODS

func (client *KrakenClient) SubscribeToTrades(s settings.BruitSettings, pairs []string) {
	if err := PubSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	client.WebSocket.SubscribeToTrades(pairs)
}

/****
	*Add func to check if already subscribed to OHLC Stream
	*Add func to get past OHLC data from rest API. Add to the candle map list
*****/
func (client *KrakenClient) SubscribeToOHLC(s settings.BruitSettings, pairs []types.Pairs, interval int) {
	var found bool = false
	for _, i := range kraken_data.GetOHLCIntervals() {
		if i == interval {
			found = true
			break
		}
	}

	if found == false {
		log.Println("Interval is not supported for Kraken Client OHLC Subscription")
		return
	}

	if err := PubSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	// add func here that makes request to rest OHLC to get past OHLC data. data should then be added to the OHLC map
	var wsPairs []string
	for _, pair := range pairs {
		_, err := client.GetOHLC(pair.Rest, interval)
		wsPairs = append(wsPairs, pair.WS)
		if err != nil {
			panic(err)
		}
		//log.Println(resp)
	}

	client.WebSocket.SubscribeToOHLC(wsPairs, interval)
}

// search through assetResp in client manager from state package. if base and quote fields match the holding and base currency, add wsname to a slice
func (client *KrakenClient) SubscribeToHoldingsOHLC(s settings.BruitSettings, interval int) {
	holdings := client.GetHoldingsWithoutStaking()
	var pairs []types.Pairs

	for _, holding := range holdings {
		for _, pair := range client.State.Client.GetAssetPairs() {
			if holding == pair.Base && strings.Join([]string{"Z", s.GetBaseCurrency()}, "") == pair.Quote {
				var p types.Pairs
				p.WS = pair.WsName
				p.Rest = pair.AltName
				pairs = append(pairs, p)

				/*resp, err := client.GetOHLC(p.Rest, interval)
				if err != nil {
					log.Println(err)
				}
				client.State.OnOHLCResponse()*/
			} else {
				log.Println("ERROR: Pair could not find match ", pair, pair.Base)
			}
		}
	}

	log.Println(pairs)

	client.SubscribeToOHLC(s, pairs, interval)
}

func (client *KrakenClient) PubDecoder(s settings.BruitSettings, OHLCch chan types.OHLCResponse, Tradech chan types.TradeResponse) {
	s.Add(1)
	defer s.Done()

	if err := PubSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocketPointer())
	client.WebSocket.GetPubSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PubJsonDecoder(message, s.GetLoggingSettings(), OHLCch, Tradech)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocketPointer())

	<-s.CtxDone()
	return
}

// ORDER BOOK SOCKET METHODS

// Subscribe to the order book.
func (client *KrakenClient) SubscribeToOrderBook(s settings.BruitSettings, pairs []string, depth int) {
	if err := BookSocketGuard(&client.WebSocket); err != nil {
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

func (client *KrakenClient) BookDecoder(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()

	if err := BookSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetBookSocketPointer())
	client.WebSocket.GetBookSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.BookJsonDecoder(message, s.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocketPointer())

	<-s.CtxDone()
	return
}

// PRIVATE SOCKET METHODS

func (client *KrakenClient) SubscribeToOpenOrders(s settings.BruitSettings, token string) {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

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

func (client *KrakenClient) CancelAll(s settings.BruitSettings, token string) {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	sub, _ := json.Marshal(&types.Subscribe{
		Event: "cancelAll",
		Token: token,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) CancelOrder(s settings.BruitSettings, token string, tradeIDs []string) {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	sub, _ := json.Marshal(&types.CancelOrder{
		Event: "cancelOrder",
		Token: token,
		Txid:  tradeIDs,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) AddOrder(s settings.BruitSettings, token string, otype string, ttype string, pair string, vol string, price string, testing bool) {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

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

func (client *KrakenClient) PrivDecoder(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()

	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPrivSocketPointer())
	client.WebSocket.GetPrivSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PrivJsonDecoder(message, s.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPrivSocketPointer())

	<-s.CtxDone()
	return
}
