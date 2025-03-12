package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
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

func (p *BackTest) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB) {
	return
}

func (p *BackTest) Stop() {
	return
}

func (p *BackTest) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	return
}
