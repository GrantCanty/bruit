package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var ErrSubscribeToHoldingsOHLC error = errors.New("No match for trading pair")
var ErrSubscribeToOHLCInterval error = errors.New("interval not supported")

func remove(slice []string, pos int) []string {
	return append(slice[:pos], slice[pos+1:]...)
}

// PUBLIC SOCKET METHODS

func (client *KrakenClient) SubscribeToTrades(s settings.BruitSettings, pairs []string) error {
	if err := PubSocketGuard(&client.WebSocket); err != nil {
		return err

	}

	if err := client.WebSocket.SubscribeToTrades(pairs); err != nil {
		return err
	}
	return nil
}

/****
	*Add func to check if already subscribed to OHLC Stream
	*Add func to get past OHLC data from rest API. Add to the candle map list
*****/
func (client *KrakenClient) SubscribeToOHLC(s settings.BruitSettings, pairs []types.Pairs, interval int) error {
	var found bool = false
	for _, i := range kraken_data.GetOHLCIntervals() {
		if i == interval {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%s - interval: %d", ErrSubscribeToOHLCInterval, interval)
	}

	if err := PubSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		return err
	}

	// add func here that makes request to rest OHLC to get past OHLC data. data should then be added to the OHLC map
	var wsPairs []string
	for _, pair := range pairs {
		_, err := client.GetOHLC(pair.Rest, interval)
		if err != nil {
			return err
		}
		wsPairs = append(wsPairs, pair.WS)
	}

	client.WebSocket.SubscribeToOHLC(wsPairs, interval)
	return nil
}

// search through assetResp in client manager from state package. if base and quote fields match the holding and base currency, add wsname to a slice
func (client *KrakenClient) SubscribeToHoldingsOHLC(s settings.BruitSettings, interval int) error {
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
				return fmt.Errorf("%s - pair: %v base: %v", ErrSubscribeToHoldingsOHLC, pair, pair.Base)
			}
		}
	}

	log.Println(pairs)

	client.SubscribeToOHLC(s, pairs, interval)
	return nil
}

func (client *KrakenClient) PubDecoder(s settings.BruitSettings, OHLCch chan types.OHLCResponse, Tradech chan types.TradeResponse, OHLCsubch chan types.OHLCSuccessResponse) {
	s.Add(1)
	defer s.Done()

	if err := PubSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocket())
	client.WebSocket.GetPubSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.PubJsonDecoder(message, s.GetLoggingSettings(), OHLCch, Tradech, OHLCsubch)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocket())

	<-s.CtxDone()
}

// ORDER BOOK SOCKET METHODS

// Subscribe to the order book.
func (client *KrakenClient) SubscribeToOrderBook(s settings.BruitSettings, depth int) {
	holdings := client.GetHoldingsWithoutStaking()

	if err := BookSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	var pairs []types.Pairs

	for _, holding := range holdings {
		for _, pair := range client.State.Client.GetAssetPairs() {
			if holding == pair.Base && strings.Join([]string{"Z", s.GetBaseCurrency()}, "") == pair.Quote {
				var p types.Pairs
				p.WS = pair.WsName
				p.Rest = pair.AltName
				pairs = append(pairs, p)
			} else {
				log.Println("ERROR: Pair could not find match ", pair, pair.Base)
			}
		}
	}
	log.Println("pairs: ", pairs)
	var wsPairs []string
	for _, pair := range pairs {
		wsPairs = append(wsPairs, pair.WS)
	}
	client.WebSocket.SubscribeToOrderBook(wsPairs, depth)
}

// need a way to save the books to a struct. on message, we read
// the struct back so we can see how to update it and then save the copy back to the struct
// then send the struct to the chan
func (client *KrakenClient) BookDecoder(s settings.BruitSettings, Bookch chan types.BookRespV2UpdateJSON, bookDepth int) {
	s.Add(1)
	defer s.Done()

	if err := BookSocketGuard(&client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetBookSocket())
	client.WebSocket.GetBookSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.BookJsonDecoder(message, s.GetLoggingSettings(), Bookch, bookDepth)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocket())

	<-s.CtxDone()
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
	client.WebSocket.GetPrivSocket().SendBinary(sub)
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
	client.WebSocket.GetPrivSocket().SendBinary(sub)
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
	client.WebSocket.GetPrivSocket().SendBinary(sub)
}

func (client *KrakenClient) PrivDecoder(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()

	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPrivSocket())
	client.WebSocket.GetPrivSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.PrivJsonDecoder(message, s.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPrivSocket())

	<-s.CtxDone()
	return
}
