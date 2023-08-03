package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

func NewBackTestEngine(parent BruitEngine) BruitEngine {
	return newProduction(parent)
}

func newBackTest(parent BruitEngine) BruitEngine {
	return &Production{BruitEngine: parent, c: nil, s: nil, db: nil}
}

type BackTest struct {
	BruitEngine

	c  clients.BruitClient
	s  settings.BruitSettings
	db *influx.DB
}

func (p *BackTest) Init(s settings.BruitSettings, c clients.BruitClient, db *influx.DB) {
	p.s = s
	p.c = c
	p.db = db

	p.s.InitSettings()
	p.c.InitClient(p.s)
	p.db.InitDB()

	return
}

func (p *BackTest) Run(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	return
}

func (p *BackTest) Stop() {
	return
}

func (p *BackTest) Wait(s settings.BruitSettings, c clients.BruitClient) {
	return
}
