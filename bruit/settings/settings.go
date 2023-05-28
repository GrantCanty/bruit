package settings

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Settings interface {
	Init()
	Wait()
	Add()
}

type DefaultSettings struct {
	ConcurrencySettings ConcurrencySettings
	GlobalSettings      Globals
	env                 map[string]string
}

type ConcurrencySettings struct {
	comms  chan os.Signal
	signal os.Signal
	Ctx    context.Context
	cancel context.CancelFunc
	Wg     sync.WaitGroup
}

type Globals struct {
	RunTime      RunTimeSettings
	Logging      LoggingSettings
	BaseCurrency string
	RiskLevel    RiskSettings
	BookDepth    int
}

type RunTimeSettings struct {
	IsSystemsTesting    bool
	IsBackTesting       bool
	IsPaperTrading      bool
	IsProductionTrading bool
}

type LoggingSettings struct {
	isLoggingToDB      bool
	isLoggingToConsole bool
}

type RiskSettings struct {
	AssetsTraded TradableAssetSettings
	Risk         RiskLevels
}

type TradableAssetSettings struct {
	TradesCrypto bool
	TradesStocks bool
	TradesForex  bool
}

type RiskLevels struct {
	isLevelOne   bool
	isLevelTwo   bool
	isLevelThree bool
}

func (s *DefaultSettings) Init() {
	s.ConcurrencySettings.Init()
	s.GlobalSettings.Init()
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

/*func (g *ConcurrencySettings) Wait() {
	<-g.comms
	g.cancel()
	g.Wg.Wait()
}*/

func (s *DefaultSettings) Wait() {
	<-s.ConcurrencySettings.comms
	s.ConcurrencySettings.cancel()
	s.ConcurrencySettings.Wg.Wait()
}

/*func (s *Settings) GetEnv() map[string]string {
	return s.env
}*/

func (g *Globals) Init() {
	g.RunTime.init()
	g.Logging.init()
	g.BaseCurrency = "USD"
	g.BookDepth = 100
	g.RiskLevel.init()
}

func (r *RunTimeSettings) init() {
	r.IsSystemsTesting = true
	r.IsBackTesting = false
	r.IsPaperTrading = false
	r.IsProductionTrading = false
}

func (l *LoggingSettings) init() {
	l.isLoggingToDB = true
	l.isLoggingToConsole = true
}

func (l LoggingSettings) GetLoggingDB() bool {
	return l.isLoggingToDB
}

func (l LoggingSettings) GetLoggingConsole() bool {
	return l.isLoggingToConsole
}

func (r *RiskSettings) init() {
	r.AssetsTraded.init()
	r.Risk.init()
}

func (a *TradableAssetSettings) init() {
	a.TradesCrypto = true
	a.TradesForex = false
	a.TradesStocks = false
}

func (l *RiskLevels) init() {
	l.isLevelOne = false
	l.isLevelTwo = false
	l.isLevelThree = true
}
