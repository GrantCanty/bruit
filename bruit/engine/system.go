package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"bruit/bruit/ws_client"
	"log"
	"os"
	"syscall"
	"time"
)

func NewSystemsTestingEngine(parent BruitEngine) BruitEngine {
	return newSystemsTesting(parent)
}

func newSystemsTesting(parent BruitEngine) BruitEngine {
	return &SystemsTesting{BruitEngine: parent}
}

type SystemsTesting struct {
	BruitEngine

	/*c  clients.BruitCryptoClient
	s  settings.BruitSettings
	db *influx.DB*/
}

func (p *SystemsTesting) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	s.Add(1)
	defer s.Done()

	krakenClient := c.(*kraken.KrakenClient)

	OHLCch := make(chan types.OHLCResponse, 1024)
	Tradech := make(chan types.TradeResponse, 1024)
	OHLCSubch := make(chan types.OHLCSuccessResponse, 16)

	ohlcMap := shared_types.OHLCVals{}

	go c.PubDecoder(s, OHLCch, Tradech, OHLCSubch)
	go func(ohlc chan types.OHLCResponse, trade chan types.TradeResponse, ohlcsub chan types.OHLCSuccessResponse, ohlcMap *shared_types.OHLCVals) {
		for {
			select {
			case res := <-ohlc:
				log.Println("ohlcResponse res: ", res)
				HandleOHLCResponse(s, c, db, res, ohlcMap)
			case res := <-trade:
				log.Println("tradeResponse res: ", res)
			case res := <-ohlcsub:
				log.Println("ohlcsub res: ", res)
				c.HandleOHLCSuccessResponse(res)
			case <-s.CtxDone():
				return
			}
		}
	}(OHLCch, Tradech, OHLCSubch, &ohlcMap)

	krakenClient.WebSocket.GetPubSocket().OnDisconnected = func(err error, socket *ws_client.Socket) {
		log.Println("disconnected from pub socket: ", err)
		go func() {
			for i := 0; ; i++ {
				// abort if the engine is shutting down
				select {
				case <-s.CtxDone():
					return
				default:
				}

				socket.Connect()
				if socket.GetIsConnected() {
					log.Println("successfully reconnected to pub socket")
					break
				}
				// cap the shift to prevent integer overflow
				shift := i
				if shift > 7 {
					shift = 7
				}

				// 32 << 7 = 4096
				backoff := time.Duration(32 << shift)
				time.Sleep(min(backoff, 4096) * time.Millisecond)
			}

			// add error handling and skip handling here
			c.SubscribeToHoldingsOHLC(s, 1)
		}()
	}

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)

	// add error handling and skip handling here
	c.SubscribeToHoldingsOHLC(s, 1)

	orderBookCh := make(chan types.BookRespV2UpdateJSON, 1024)
	orderBookErrCh := make(chan error, 1)
	var bookDepth int = 10

	go c.BookDecoder(s, orderBookCh, orderBookErrCh, bookDepth)
	go func(book chan types.BookRespV2UpdateJSON, errCh chan error) {
		for {
			select {
			case res := <-book:
				log.Println("orderbook: ", res)
			case err := <-errCh:
				log.Println("orderbook error: ", err)

				if IsFatalErr(err) {
					p.Stop()
					return
				}

			case <-s.CtxDone():
				return
			}
		}
	}(orderBookCh, orderBookErrCh)

	krakenClient.WebSocket.GetBookSocket().OnDisconnected = func(err error, socket *ws_client.Socket) {
		log.Println("disconnected from book socket: ", err)
		go func() {
			for i := 0; ; i++ {
				// abort if the engine is shutting down
				select {
				case <-s.CtxDone():
					return
				default:
				}

				socket.Connect()
				if socket.GetIsConnected() {
					log.Println("successfully reconnected to book socket")
					break
				}
				// cap the shift to prevent integer overflow
				shift := i
				if shift > 7 {
					shift = 7
				}

				// 32 << 7 = 4096
				backoff := time.Duration(32 << shift)
				time.Sleep(min(backoff, 4096) * time.Millisecond)
			}

			c.SubscribeToOrderBook(s, bookDepth)
		}()
	}

	c.SubscribeToOrderBook(s, bookDepth)

	<-s.CtxDone()
}

func (p *SystemsTesting) Stop() {
	log.Println("FATAL ERROR: Shutting down systems testing engine gracefully...")
	if pid := os.Getpid(); pid > 0 {
		syscall.Kill(pid, syscall.SIGINT)
	}
}

func (p *SystemsTesting) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}

func HandleOHLCResponse(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, data types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	/**
	*  Add:
	*  OHLCResponseHandler func to add responses to a LL. should delete the head if length is too long (ex: 10,000)
	*  CalcTechnicals func to recalculate the values of technical indicators
	*  Eval func to evaluate if buy/sell condition is met
	*  PlaceOrder func depending on Eval func
	**/
	c.HandleOHLCResponse(data, ohlcMap)
}
