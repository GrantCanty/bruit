package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"log"
	"time"
)

func NewPaperTradingEngine(parent BruitEngine) BruitEngine {
	return newPaperTrading(parent)
}

func newPaperTrading(parent BruitEngine) BruitEngine {
	//return &Production{BruitEngine: parent}
	return &PaperTrading{BruitEngine: parent}
}

type PaperTrading struct {
	BruitEngine

	//c  clients.BruitCryptoClient
	//s  settings.BruitSettings
	//db *influx.DB
}

func (p *PaperTrading) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	s.Add(1)
	defer s.Done()

	OHLCch := make(chan types.OHLCResponse, 1024)
	Tradech := make(chan types.TradeResponse, 1024)
	// OHLCSubch := make(chan types.OHLCSuccessResponse, 16)

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
			case <-s.CtxDone():
				return
			}

		}
	}(OHLCch, Tradech)

	<-s.CtxDone()
}

func (p *PaperTrading) Stop() {
	return
}

func (p *PaperTrading) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
