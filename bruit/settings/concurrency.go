package settings

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ConcurrencySettings struct {
	comms  chan os.Signal
	signal os.Signal
	Ctx    context.Context
	cancel context.CancelFunc
	Wg     sync.WaitGroup
}

func (g *ConcurrencySettings) ReturnWg() *sync.WaitGroup {
	return &g.Wg
}

func (g *ConcurrencySettings) Init() {
	g.comms = make(chan os.Signal, 1)
	signal.Notify(g.comms, os.Interrupt, syscall.SIGTERM)
	g.Ctx = context.Background()
	g.Ctx, g.cancel = context.WithCancel(g.Ctx)
}
