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
	//return &Production{BruitEngine: parent}
	return &Production{BruitEngine: parent, c: nil, s: nil, db: nil}
}

type Production struct {
	BruitEngine

	c  clients.BruitClient
	s  settings.BruitSettings
	db *influx.DB
}

func (p *Production) Init(s settings.BruitSettings, c clients.BruitClient, db *influx.DB) {
	p.s = s
	p.c = c
	p.db = db

	p.s.InitSettings()
	p.c.InitClient(s)
	p.db.InitDB()
}

func (p *Production) Run( /*s settings.BruitSettings, c clients.BruitClient, db influx.DB*/ ) {
	p.s.Add(1)
	defer p.s.Done()

	go p.c.PubDecoder(p.s)

	ohlcMap := shared_types.OHLCVals{}
	go p.c.PubListen(p.s, &ohlcMap, p.db.GetTradeWriter())

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	p.c.SubscribeToHoldingsOHLC(p.s, 1)

	<-p.s.CtxDone()
}

func (p *Production) Stop() {
	return
}

func (p *Production) Wait( /*s settings.BruitSettings, c clients.BruitClient*/ ) {
	go p.c.DeferChanClose(p.s)
	p.s.Wait()
}
