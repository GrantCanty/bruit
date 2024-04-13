package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
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

func (p *PaperTrading) Init(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, str shared_types.Strategy) {
	//p.s = s
	//p.c = c
	//p.db = db

	//p.s.InitSettings()
	//p.c.InitClient(s)
	//p.db.InitDB()
}

func (p *PaperTrading) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	s.Add(1)
	defer s.Done()

	go c.PubDecoder(s)

	//ohlcMap := shared_types.OHLCVals{}
	//go c.PubListen(s, &ohlcMap, db.GetTradeWriter())

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)

	<-s.CtxDone()
}

func (p *PaperTrading) Stop() {
	return
}

func (p *PaperTrading) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
