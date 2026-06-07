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
	"time"
)

var ErrPairNotFound error = errors.New("No match for trading pair")
var ErrSubscribeToOHLCInterval error = errors.New("Interval not supported")

var ErrFailedToMarshalCancelAll error = errors.New("Failed to marshal CancelAll message")
var ErrFailedToMarshalCancelOrder error = errors.New("Failed to marshal CancelOrder message")
var ErrFailedToMarshalAddOrder error = errors.New("Failed to marshal AddOrder message")

var ErrFailedToSendCancelAll error = errors.New("Failed to send CancelAll message")
var ErrFailedToSendCancelOrder error = errors.New("Failed to send CancelOrder message")
var ErrFailedToSendAddOrder error = errors.New("Failed to send AddOrder message")

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

	return client.WebSocket.SubscribeToOHLC(wsPairs, interval)
}

// search through assetResp in client manager from state package. if base and quote fields match the holding and base currency, add wsname to a slice
func (client *KrakenClient) SubscribeToHoldingsOHLC(s settings.BruitSettings, interval int) (skipped []string, err error) {
	holdings := client.GetHoldingsWithoutStaking()
	var pairs []types.Pairs

	for _, holding := range holdings {
		var found bool
		for _, pair := range client.State.Client.GetAssetPairs() {
			if holding == pair.Base && strings.Join([]string{"Z", s.GetBaseCurrency()}, "") == pair.Quote {
				var p types.Pairs
				p.WS = pair.WsName
				p.Rest = pair.AltName
				pairs = append(pairs, p)
				found = true
				break
			}
		}
		if !found {
			skipped = append(skipped, holding)
			log.Printf("Warning: no trading pair found for holding %s with base currency %s (skipping)", holding, s.GetBaseCurrency())
		}
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("%w - zero trading pairs matched for holdings: %v", ErrPairNotFound, holdings)
	}

	if err := client.SubscribeToOHLC(s, pairs, interval); err != nil {
		return nil, err
	}
	return skipped, nil
}

func (client *KrakenClient) PubDecoder(s settings.BruitSettings, OHLCch chan types.OHLCResponse, Tradech chan types.TradeResponse, OHLCsubch chan types.OHLCSuccessResponse) {
	s.Add(1)
	defer s.Done()

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocket())
	client.WebSocket.GetPubSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.PubJsonDecoder(message, s.GetLoggingSettings(), OHLCch, Tradech, OHLCsubch)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocket())

	var err error
	maxRetries := 5
	backoff := 32 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err = PubSocketGuard(&client.WebSocket)
		if err == nil {
			break
		}

		if errors.Is(err, ErrPubSocketNotInit) || errors.Is(err, ErrNotPubSocket) {
			log.Fatalf("FATAL - Dev configuration error in PubDecoder: %v", err)
		}

		log.Printf("Warning: Connection attempt %d failed: %v. Retrying in %v...", i+1, err, backoff)

		select {
		case <-s.CtxDone():
			return
		case <-time.After(backoff):
		}

		backoff *= 2

	}

	if err != nil {
		log.Printf("ERROR - PubDecoder failed to establish a connection to the Kraken WebSocket server after %d attempts: %v\n", maxRetries, err)
		return
	}

	<-s.CtxDone()
}

// ORDER BOOK SOCKET METHODS

// Subscribe to the order book.
func (client *KrakenClient) SubscribeToOrderBook(s settings.BruitSettings, depth int) (skipped []string, err error) {
	holdings := client.GetHoldingsWithoutStaking()

	if err := BookSocketGuard(&client.WebSocket); err != nil {
		return nil, err
	}

	var pairs []types.Pairs

	for _, holding := range holdings {
		var found bool
		for _, pair := range client.State.Client.GetAssetPairs() {
			if holding == pair.Base && strings.Join([]string{"Z", s.GetBaseCurrency()}, "") == pair.Quote {
				var p types.Pairs
				p.WS = pair.WsName
				p.Rest = pair.AltName
				pairs = append(pairs, p)
				found = true
				break
			}
		}
		if !found {
			skipped = append(skipped, holding)
			log.Printf("Warning - no order book pair found for holding %s with base currency %s (skipping)", holding, s.GetBaseCurrency())
		}
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("%w - zero order book pairs matched for holdings: %v", ErrPairNotFound, holdings)
	}

	var wsPairs []string
	for _, pair := range pairs {
		wsPairs = append(wsPairs, pair.WS)
	}
	if err := client.WebSocket.SubscribeToOrderBook(wsPairs, depth); err != nil {
		return nil, err
	}
	return skipped, nil
}

// need a way to save the books to a struct. on message, we read
// the struct back so we can see how to update it and then save the copy back to the struct
// then send the struct to the chan
func (client *KrakenClient) BookDecoder(s settings.BruitSettings, Bookch chan types.BookRespV2UpdateJSON, bookDepth int) {
	s.Add(1)
	defer s.Done()

	ws_client.ReceiveLocker(client.WebSocket.GetBookSocket())
	client.WebSocket.GetBookSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.BookJsonDecoder(message, s.GetLoggingSettings(), Bookch, bookDepth)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocket())

	var err error
	maxRetries := 5
	backoff := 32 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err = BookSocketGuard(&client.WebSocket)
		if err == nil {
			break
		}

		if errors.Is(err, ErrBookSocketNotInit) || errors.Is(err, ErrNotBookSocket) {
			log.Fatalf("FATAL - Dev configuration error in BookDecoder: %v", err)
		}

		log.Printf("Warning: Connection attempt %d failed: %v. Retrying in %v...", i+1, err, backoff)

		select {
		case <-s.CtxDone():
			return
		case <-time.After(backoff):
		}

		backoff *= 2

	}

	if err != nil {
		log.Printf("ERROR - BookDecoder failed to establish a connection to the Kraken WebSocket server after %d attempts: %v\n", maxRetries, err)
		return
	}

	<-s.CtxDone()
}

// PRIVATE SOCKET METHODS

func (client *KrakenClient) SubscribeToOpenOrders(s settings.BruitSettings, token string) error {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		return err
	}

	return client.WebSocket.SubscribeToOpenOrders(token)
}

func (client *KrakenClient) CancelAll(s settings.BruitSettings, token string) error {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		return err
	}

	sub, err := json.Marshal(&types.Subscribe{
		Event: "cancelAll",
		Token: token,
	})

	if err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToMarshalCancelAll, err)
	}

	if err := client.WebSocket.GetPrivSocket().SendBinary(sub); err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToSendCancelAll, err)
	}
	return nil
}

func (client *KrakenClient) CancelOrder(s settings.BruitSettings, token string, tradeIDs []string) error {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		return err
	}

	sub, err := json.Marshal(&types.CancelOrder{
		Event: "cancelOrder",
		Token: token,
		Txid:  tradeIDs,
	})

	if err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToMarshalCancelOrder, err)
	}

	if err := client.WebSocket.GetPrivSocket().SendBinary(sub); err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToSendCancelOrder, err)
	}
	return nil
}

func (client *KrakenClient) AddOrder(s settings.BruitSettings, token string, otype string, ttype string, pair string, vol string, price string, testing bool) error {
	if err := PrivSocketGuard(&client.WebSocket); err != nil {
		return err
	}

	test := strconv.FormatBool(testing)
	sub, err := json.Marshal(&types.Order{
		WsToken:   token,
		Event:     "addOrder",
		OrderType: otype,
		TradeType: ttype,
		Pair:      pair,
		Volume:    vol,
		Price:     price,
		Validate:  test,
	})

	if err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToMarshalAddOrder, err)
	}

	if err := client.WebSocket.GetPrivSocket().SendBinary(sub); err != nil {
		return fmt.Errorf("%s - %w", ErrFailedToSendAddOrder, err)
	}
	return nil
}

func (client *KrakenClient) PrivDecoder(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()

	ws_client.ReceiveLocker(client.WebSocket.GetPrivSocket())
	client.WebSocket.GetPrivSocket().OnTextMessage = func(message []byte, socket *ws_client.Socket) {
		client.WebSocket.PrivJsonDecoder(message, s.GetLoggingSettings())
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPrivSocket())

	var err error
	maxRetries := 5
	backoff := 32 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err = PrivSocketGuard(&client.WebSocket)
		if err == nil {
			break
		}

		if errors.Is(err, ErrPrivSocketNotInit) || errors.Is(err, ErrNotPrivSocket) {
			log.Fatalf("FATAL - Dev configuration error in PrivDecoder: %v", err)
		}

		log.Printf("Warning: Connection attempt %d failed: %v. Retrying in %v...", i+1, err, backoff)

		select {
		case <-s.CtxDone():
			return
		case <-time.After(backoff):
		}

		backoff *= 2

	}

	if err != nil {
		log.Printf("ERROR - PrivDecoder failed to establish a connection to the Kraken WebSocket server after %d attempts: %v\n", maxRetries, err)
		return
	}

	<-s.CtxDone()
}
