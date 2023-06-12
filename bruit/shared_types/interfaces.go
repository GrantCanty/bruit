package shared_types

import (
	bruit "bruit/bruit/settings"
	"time"

	//"bruit/bruit/clients/kraken/types"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shopspring/decimal"
)

type WebSocketClient interface {
	// GENERAL METHODS
	InitWebSockets(g *bruit.BruitSettings)
	DeferChanClose(g *bruit.BruitSettings)

	// PUBLIC SOCKET METHODS
	SubscribeToOHLC(g *bruit.BruitSettings, pairs []string, interval int)
	SubscribeToTrades(g *bruit.BruitSettings, pairs []string)
	PubDecoder(g *bruit.BruitSettings)
	PubListen(g *bruit.BruitSettings, ohlcMap *OHLCVals, tradesWriter api.WriteAPI) // needs to take an interface instead of OHLCVals

	// ORDER BOOK SOCKET METHODS
	SubscribeToOrderBook(g *bruit.BruitSettings, pairs []string, depth int)
	//BookDecode(g *ConcurrencySettings)
	BookListen(g *bruit.BruitSettings)

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
	SetClose(num decimal.Decimal, vol decimal.Decimal)
	GetVWAP() decimal.Decimal
	SetVWAP(num decimal.Decimal, vol decimal.Decimal)
	GetVolume() decimal.Decimal
	SetVolume(num decimal.Decimal)
	GetCount() int
	SetCount(num int, vol decimal.Decimal)
}

/*type OHLCValHolder interface {
	Set(key int, data List)
	RLock()
	RUnlock()
	Lock()
	Unlock()
	GetData() map[SubscriptionMetaData]List
	GetMutex() *sync.RWMutex
}*/

type SubscriptionMetaData interface {
	GetChannelID() int
	GetChannelName() string
	GetPair() string
	Found(metaData SubscriptionMetaData) bool
}

type SubscriptionData interface {
	GetData() SubscriptionData
}
