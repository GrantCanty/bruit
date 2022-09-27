package shared_types

import (
	"bruit/bruit"
	"sync"
	"time"

	//"bruit/bruit/clients/kraken/types"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shopspring/decimal"
)

type WebSocketClient interface {
	// GENERAL METHODS
	InitWebSockets(g *bruit.Settings)
	DeferChanClose(g *bruit.Settings)

	// PUBLIC SOCKET METHODS
	SubscribeToOHLC(g *bruit.Settings, pairs []string, interval int)
	SubscribeToTrades(g *bruit.Settings, pairs []string)
	PubDecoder(g *bruit.Settings)
	PubListen(g *bruit.Settings, ohlcMap OHLCValHolder, tradesWriter api.WriteAPI) // needs to take an interface instead of OHLCVals

	// ORDER BOOK SOCKET METHODS
	SubscribeToOrderBook(g *bruit.Settings, pairs []string, depth int)
	//BookDecode(g *ConcurrencySettings)
	BookListen(g *bruit.Settings)

	// PRIVATE SOCKET METHODS

	//PrivListen()
}

type Candle interface {
	SetCandle(data ...interface{})
	GetCandle() Candle
	GetStartTime() UnixTime
	SetStartTime(newTime time.Time)
	GetEndTime() UnixTime
	SetEndTime(newTime time.Time)
	GetHigh() decimal.Decimal
	SetHigh(num decimal.Decimal)
	GetLow() decimal.Decimal
	SetLow(num decimal.Decimal)
	GetClose() decimal.Decimal
	SetClose(num decimal.Decimal, num2 decimal.Decimal)
	GetVWAP() decimal.Decimal
	SetVWAP(num decimal.Decimal, num2 decimal.Decimal)
	GetVolume() decimal.Decimal
	SetVolume(num decimal.Decimal)
	GetCount() int
	SetCount(num int, num2 decimal.Decimal)
}

type OhlcResponseHolder interface {
	GetChannelID() int
	GetList() *List
	GetInterval() int64
	//Return() *OhlcResponseHolder
}

type OHLCValHolder interface {
	Set(key int, data OhlcResponseHolder)
	RLock()
	RUnlock()
	Lock()
	Unlock()
	GetData() map[int]OhlcResponseHolder
	GetMutex() *sync.RWMutex
}
