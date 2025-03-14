package web_socket

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	decoders "bruit/bruit/clients/kraken/decoder_funcs"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type WebSocketClient struct {
	pubSocket  ws_client.Socket
	bookSocket ws_client.Socket
	privSocket ws_client.Socket

	bookChan     chan interface{} // this chan contains the final book data
	bookJSONChan chan interface{} // this chan contains the most recently decoded data from the book subscription
	privChan     chan interface{}

	// switch to ordered map instead
	orderBooks      map[string]*types.OrderBookWithMutexTree
	orderBooksMutex sync.RWMutex
}

type MessageTypeIdentifier struct {
	Channel string `json:"channel"`
	Type    string `json:"type"`
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

func (client *WebSocketClient) BookJsonDecoder(response string, logger settings.LoggingSettings, Bookch chan types.BookRespV2UpdateJSON) {
	byteResponse := []byte(response)

	var msgType MessageTypeIdentifier
	if err := json.Unmarshal(byteResponse, &msgType); err != nil {
		log.Println("Error identifying message type:", err)
		return
	}

	switch msgType.Channel {
	case "book":
		switch msgType.Type {
		case "update":
			if resp, err := decoders.UpdateBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
				symbol := resp.Data[0].Symbol

				client.orderBooksMutex.Lock()
				client.orderBooks[symbol].Mutex.Lock()
				book := client.orderBooks[symbol].Book

				for _, bid := range resp.Data[0].Bids {
					if val, err := bid.Quantity.Float64(); err == nil {
						if val == 0 {
							if _, ok := book.Bids.Get(bid.Price); ok {
								book.Bids.Remove(bid.Price)
							} else {
								panic("bid not found. book out of order")
							}
						} else {
							book.Bids.Put(bid.Price, bid.Quantity)
						}
					}
				}

				for _, ask := range resp.Data[0].Asks {
					if val, err := ask.Quantity.Float64(); err == nil {
						if val == 0 {
							if _, ok := book.Asks.Get(ask.Price); ok {
								book.Asks.Remove(ask.Price)
							} else {
								panic("ask not found. book out of order")
							}
						} else {
							book.Asks.Put(ask.Price, ask.Quantity)
						}
					}
				}

				if book.Bids.Size() > 10 {
					keys := book.Bids.Keys()

					for i := 10; i < len(keys); i++ {
						book.Bids.Remove(keys[i])
					}
				}

				if book.Asks.Size() > 10 {
					keys := book.Asks.Keys()

					for i := 10; i < len(keys); i++ {
						book.Asks.Remove(keys[i])
					}
				}

				if ok := decoders.VerifyChecksumUpdate(*book, *resp); !ok {
					panic("checksums don't match")
				}

				client.orderBooksMutex.Unlock()
				client.orderBooks[symbol].Mutex.Unlock()
				Bookch <- *book
			}

		case "snapshot":
			t := time.Now()
			if resp, err := decoders.SnapshotBookResponseDecoderV2(byteResponse, logger.GetLoggingConsole()); err == nil {
				symbol := resp.Data[0].Symbol
				log.Println("symbol: ", symbol, t)

				bidsMap := treemap.NewWith(func(key interface{}, value interface{}) int {
					return -types.NumericStringComparator(key, value) // Descending order
				})

				for _, bid := range resp.Data[0].Bids {
					bidsMap.Put(bid.Price, bid.Quantity)
				}

				asksMap := treemap.NewWith(types.NumericStringComparator)
				for _, ask := range resp.Data[0].Asks {
					asksMap.Put(ask.Price, ask.Quantity)
				}

				book := &types.OrderBookWithMutexTree{
					Book: &types.BookRespV2UpdateJSON{
						BookRespV2SnapshotJSON: types.BookRespV2SnapshotJSON{
							Symbol:   symbol,
							Asks:     asksMap,
							Bids:     bidsMap,
							Checksum: resp.Data[0].Checksum,
						},
						Timestamp: t,
					},
					Mutex: sync.RWMutex{},
				}

				client.orderBooksMutex.Lock()
				client.orderBooks[symbol] = book
				client.orderBooksMutex.Unlock()
				Bookch <- *book.Book
			}

		}
	}

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
	ws.bookChan = make(chan interface{})
	ws.bookJSONChan = make(chan interface{})
	ws.privChan = make(chan interface{})
}

func (ws *WebSocketClient) InitBook() {
	ws.orderBooks = make(map[string]*types.OrderBookWithMutexTree)
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
			Depth:   depth,
			Channel: "book",
			Symbol:  pairs,
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
