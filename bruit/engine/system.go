package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"log"

	//"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	//"log"
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

	//var OHLCch chan types.OHLCResponse
	/*OHLCch := make(chan types.OHLCResponse)

	//var Tradech chan types.TradeResponse
	Tradech := make(chan types.TradeResponse)

	//var OHLCSubch chan types.OHLCSuccessResponse
	OHLCSubch := make(chan types.OHLCSuccessResponse)

	go c.PubDecoder(s, OHLCch, Tradech, OHLCSubch)

	go func(ohlc chan types.OHLCResponse, trade chan types.TradeResponse, ohlcsub chan types.OHLCSuccessResponse) {
		for {
			select {
			case res := <-ohlc:
				log.Println("ohlcResponse res: ", res)
			case res := <-trade:
				log.Println("tradeResponse res: ", res)
			case res := <-ohlcsub:
				log.Println("ohlcsub res: ", res)
				c.HandleOHLCSuccessResponse(res)
			}
		}
	}(OHLCch, Tradech, OHLCSubch)

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)*/

	orderBookCh := make(chan types.BookRespV2Update)

	go c.BookDecoder(s, orderBookCh)
	go func(book chan types.BookRespV2Update) {
		for {
			select {
			case res := <-book:
				log.Println(res)
			}
		}
	}(orderBookCh)

	c.SubscribeToOrderBook(s, 10)

	<-s.CtxDone()
}

func (p *SystemsTesting) Stop() {
	return
}

func (p *SystemsTesting) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
