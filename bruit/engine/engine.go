package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
)

type BruitEngine interface {
	Init(s settings.BruitSettings, c clients.BruitClient, db influx.DB)
	Run(s settings.BruitSettings, c clients.BruitClient, db influx.DB)
	Stop()
	Wait(s settings.BruitSettings, c clients.BruitClient)
}

type emptyEngine int

func (e emptyEngine) Init(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Run(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Stop() {
	return
}

func (e emptyEngine) Wait(s settings.BruitSettings, c clients.BruitClient) {
	return
}

func Engine() BruitEngine {
	return new(emptyEngine)
}

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

	c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)

	<-s.CtxDone()
}

func (p *Production) Stop() {
	return
}

func (p *Production) Wait(s settings.BruitSettings, c clients.BruitClient) {
	//go p.c.DeferChanClose(s)
	go c.DeferChanClose(s)
	s.Wait()
}
