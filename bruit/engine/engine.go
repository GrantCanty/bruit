package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"log"
	"time"
)

type BruitEngine interface {
	Init(s settings.Settings, c clients.BruitClient, db influx.DB)
	Run(s settings.Settings)
	Stop()
	Wait(s settings.Settings, c clients.BruitClient)
}

type emptyEngine int

func (e emptyEngine) Init(s settings.Settings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Run(s settings.Settings) {
	return
}

func (e emptyEngine) Stop() {
	return
}

func (e emptyEngine) Wait(s settings.Settings, c clients.BruitClient) {
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

func (p *Production) Init(s settings.Settings, c clients.BruitClient, db influx.DB) {
	s.Init()
	c.InitClient(s)
	db.Init()
}

func (p *Production) Run(s settings.Settings) {
	s.Add(1)
	defer s.Done()

	timer := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case t := <-timer.C:
				log.Println("tick at", t)
				/*case <-s.CtxDone():
				log.Println("stopping")
				return*/
			}

		}
	}()

	<-s.CtxDone()
}

func (p *Production) Stop() {
	return
}

func (p *Production) Wait(s settings.Settings, c clients.BruitClient) {
	//go p.c.DeferChanClose(s)
	go c.DeferChanClose(s)
	s.Wait()
}
