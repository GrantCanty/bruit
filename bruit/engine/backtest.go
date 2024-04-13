package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"log"
)

func NewBackTestEngine(parent BruitEngine) BruitEngine {
	return newProduction(parent)
}

func newBackTest(parent BruitEngine) BruitEngine {
	return &BackTest{BruitEngine: parent}
}

type BackTest struct {
	BruitEngine

	/*c  clients.BruitCryptoClient
	s  settings.BruitSettings
	db *influx.DB*/
	ohlcData types.OHLCResp
}

func (p *BackTest) Init(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, str shared_types.Strategy) {
	/*p.s = s
	p.c = c
	p.db = db*/

	/*s.InitSettings()
	c.InitClient(s)
	p.db.InitDB()*/

	ohlcData, err := c.GetOHLC("DOT/USD", 5)
	if err != nil {
		log.Println(err)
	}
	log.Println(ohlcData)

	return
}

func (p *BackTest) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	return
}

func (p *BackTest) Stop() {
	return
}

func (p *BackTest) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	return
}
