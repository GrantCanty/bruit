package engine

import (
	"bruit/bruit"
	"bruit/bruit/clients"
	"bruit/bruit/influx"
)

type BruitEngine interface {
	Init(s *bruit.Settings, c clients.BruitClient, db influx.DB)
	Run()
	Stop()
}

type emptyEngine int

func (e emptyEngine) Init(s *bruit.Settings, c clients.BruitClient, db influx.DB) {
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

func ProductionEngine(parent BruitEngine) (engine BruitEngine) {
	return newProduction(parent)
}

func newProduction(parent BruitEngine) BruitEngine {
	return &Production{BruitEngine: parent}
}

type Production struct {
	BruitEngine
}

func (p *Production) Init(s *bruit.Settings, c clients.BruitClient, db influx.DB) {
	s.Init()
	c.InitClient(s)
	db.Init()
}

func (p *Production) Run()

func (p *Production) Stop()
