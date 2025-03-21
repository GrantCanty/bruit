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

	privChan chan interface{}

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
}

func (client *WebSocketClient) BookJsonDecoder(response string, logger settings.LoggingSettings, Bookch chan types.BookRespV2UpdateJSON, bookDepth int) {
	byteResponse := []byte(response)

	var msgType MessageTypeIdentifier
	if err := json.Unmarshal(byteResponse, &msgType); err != nil {
		log.Println("Error identifying message type:", err)
		return
	} else {
		log.Println("got message: ", msgType, err)
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

				if book.Bids.Size() > bookDepth {
					keys := book.Bids.Keys()

					for i := bookDepth; i < len(keys); i++ {
						book.Bids.Remove(keys[i])
					}
				}

				if book.Asks.Size() > bookDepth {
					keys := book.Asks.Keys()

					for i := bookDepth; i < len(keys); i++ {
						book.Asks.Remove(keys[i])
					}
				}

				if ok := types.VerifyChecksumUpdate(*book, *resp); !ok {
					panic("checksums don't match")
				}

				bookCopy := types.DeepCopyOrderBook(*book)
				client.orderBooksMutex.Unlock()
				client.orderBooks[symbol].Mutex.Unlock()
				Bookch <- bookCopy
				break
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
				break
			}
		default:
			log.Println("default 2. unsuccessful attempt at unmarshalling data ", response)
		}
	case "status":
		switch msgType.Type {
		case "update":
			if resp, err := decoders.StatusBookResponseV2WS(byteResponse, logger.GetLoggingConsole()); err == nil {
				log.Println("StatusBookResponseV2WS resp: ", resp)
			} else {
				log.Println("error in StatusBookResponseV2WS switch: ", err)
			}
		default:
			log.Println("in update part of switch. unknown data response: ", response)
		}
	default:
		if resp, err := decoders.SubscribeResponseV2WS(byteResponse, logger.GetLoggingConsole()); err == nil {
			log.Println("SubscribeResponseV2WS resp: ", resp)
		} else {
			log.Println("default 1. unsuccessful attempt at unmarshalling data ", response, err)
		}
	}
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

	client.privChan <- resp
}

func (ws *WebSocketClient) InitChannels() {
	ws.privChan = make(chan interface{})
}

func (ws *WebSocketClient) InitBook() {
	ws.orderBooks = make(map[string]*types.OrderBookWithMutexTree)
}

func (ws *WebSocketClient) SubscribeToTrades(pairs []string) {
	sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name: "trade",
		},
		Pair: pairs,
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
	ws.pubSocket.SendBinary(sub)
}

func (ws *WebSocketClient) SubscribeToOHLC(pairs []string, interval int) {
	sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.OHLCSubscription{
			Interval: interval,
			Name:     "ohlc",
		},
		Pair: pairs,
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
	ws.pubSocket.SendBinary(sub)
}

func (ws *WebSocketClient) SubscribeToOpenOrders(token string) {
	sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name:  "openOrders",
			Token: token,
		},
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
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

func (ws *WebSocketClient) GetPrivChan() chan interface{} {
	return ws.privChan
}
