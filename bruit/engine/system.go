package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
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

	OHLCch := make(chan types.OHLCResponse)
	Tradech := make(chan types.TradeResponse)
	OHLCSubch := make(chan types.OHLCSuccessResponse)

	ohlcMap := shared_types.OHLCVals{}

	go c.PubDecoder(s, OHLCch, Tradech, OHLCSubch)
	go func(ohlc chan types.OHLCResponse, trade chan types.TradeResponse, ohlcsub chan types.OHLCSuccessResponse, ohlcMap *shared_types.OHLCVals) {
		for {
			select {
			case res := <-ohlc:
				log.Println("ohlcResponse res: ", res)
				c.HandleOHLCResponse(res, ohlcMap)
			case res := <-trade:
				log.Println("tradeResponse res: ", res)
			case res := <-ohlcsub:
				log.Println("ohlcsub res: ", res)
				c.HandleOHLCSuccessResponse(res)
			}
		}
	}(OHLCch, Tradech, OHLCSubch, &ohlcMap)

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)

	orderBookCh := make(chan types.BookRespV2UpdateJSON)
	var bookDepth int = 10

	go c.BookDecoder(s, orderBookCh, bookDepth)
	go func(book chan types.BookRespV2UpdateJSON) {
		for res := range book {
			log.Println("orderbook: ", res)
		}
	}(orderBookCh)

	c.SubscribeToOrderBook(s, bookDepth)

	<-s.CtxDone()
}

func (p *SystemsTesting) Stop() {
	return
}

func (p *SystemsTesting) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
