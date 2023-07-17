package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
)

func NewProductionEngine(parent BruitEngine) BruitEngine {
	return newProduction(parent)
}

func newProduction(parent BruitEngine) BruitEngine {
	return &Production{BruitEngine: parent}
}

type Production struct {
	BruitEngine

	c clients.BruitClient
}

func (p *Production) Init(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	s.Init()
	c.InitClient(s)
	db.Init()
}

func (p *Production) Run(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	s.Add(1)
	defer s.Done()

	go c.PubDecoder(s)

	ohlcMap := shared_types.OHLCVals{}
	go c.PubListen(s, &ohlcMap, db.GetTradeWriter())

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)

	<-s.CtxDone()
}

func (p *Production) Stop() {
	return
}

func (p *Production) Wait(s settings.BruitSettings, c clients.BruitClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
