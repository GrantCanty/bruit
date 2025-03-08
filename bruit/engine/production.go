package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"log"
	"time"

)

func NewProductionEngine(parent BruitEngine) BruitEngine {
	return newProduction(parent)
}

func newProduction(parent BruitEngine) BruitEngine {
	//return &Production{BruitEngine: parent}
	return &Production{BruitEngine: parent}
}

type Production struct {
	BruitEngine

	//c  clients.BruitCryptoClient
	//s  settings.BruitSettings
	//db *influx.DB
}

func (p *Production) Init(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, str shared_types.Strategy) {
	//p.s = s
	//p.c = c
	//p.db = db

	//p.s.InitSettings()
	//p.c.InitClient(s)
	//p.db.InitDB()
}

func (p *Production) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	s.Add(1)
	defer s.Done()

	var OHLCch chan types.OHLCResponse
	OHLCch = make(chan types.OHLCResponse)

	var Tradech chan types.TradeResponse
	Tradech = make(chan types.TradeResponse)

	//go c.PubDecoder(s, OHLCch, Tradech)

	//ohlcMap := shared_types.OHLCVals{}
	//go c.PubListen(s, OHLCch, Tradech)

	go func(ohlc chan types.OHLCResponse, trade chan types.TradeResponse) {
		for {
			select {
				case res := <-ohlc:
					log.Println("received response in Run function: ", time.Now())
					log.Println("ohlcResponse res: ", res)
					
				case res := <-trade:
					log.Println("tradeResponse res: ", res)
				}

		}
	}(OHLCch, Tradech)

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)

	<-s.CtxDone()
}

func (p *Production) Stop() {
	return
}

func (p *Production) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
