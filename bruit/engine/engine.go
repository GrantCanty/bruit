package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

type BruitEngine interface {
	Init(s settings.Settings, c clients.BruitClient, db influx.DB)
	Run()
	Stop()
	Wait()
}

type emptyEngine int

func (e emptyEngine) Init(s settings.Settings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Run() {
	return
}

func (e emptyEngine) Stop() {
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
}

func (p *Production) Init(s settings.Settings, c clients.BruitClient, db influx.DB) {
	s.Init()
	c.InitClient(s)
	db.Init()
}

func (p *Production) Run()

func (p *Production) Stop()
