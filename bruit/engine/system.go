package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"log"
)

func NewSystemsTestingEngine(parent BruitEngine) BruitEngine {
	return newSystemsTesting(parent)
}

func newSystemsTesting(parent BruitEngine) BruitEngine {
	//return &Production{BruitEngine: parent}
	return &SystemsTesting{BruitEngine: parent}
}

type SystemsTesting struct {
	BruitEngine

	/*c  clients.BruitCryptoClient
	s  settings.BruitSettings
	db *influx.DB*/
}

func (p *SystemsTesting) Init(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, str shared_types.Strategy) {
	/*p.s = s
	p.c = c
	p.db = db

	p.s.InitSettings()
	p.c.InitClient(s)
	p.db.InitDB()*/
}

func (p *SystemsTesting) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	s.Add(1)
	defer s.Done()

	var ch chan types.OHLCResponse
	ch = make(chan types.OHLCResponse)

	go c.PubDecoder(s)

	//ohlcMap := shared_types.OHLCVals{}
	go c.PubListen(s, ch, db.GetTradeWriter())

	go func(ch chan types.OHLCResponse) {
		for {
			select {
			case res := <-ch:
				log.Println("resssszszszs: ", res)
			}
		}
	}(ch)

	//c.SubscribeToOHLC(s, []string{"EOS/USD", "BTC/USD"}, 1)
	c.SubscribeToHoldingsOHLC(s, 1)

	<-s.CtxDone()
}

func (p *SystemsTesting) Stop() {
	return
}

func (p *SystemsTesting) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	go c.DeferChanClose(s)
	s.Wait()
}
