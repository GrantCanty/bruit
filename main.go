package main

import (
	"bruit/bruit/clients"
	"bruit/bruit/clients/kraken"
	"bruit/bruit/engine"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

func main() {
	var s settings.BruitSettings
	s = settings.NewDefaultSettings(s)

	db := influx.DB{}

	var c clients.BruitClient
	c = &kraken.KrakenClient{}

	var e engine.BruitEngine
	e = engine.NewProductionEngine(e)
	e.Init(s, c, db)
	c.GetHoldingsWithoutStaking()
	//go e.Run(s, c, db)

	e.Wait(s, c)
}
