package shared_types

import (
	"bruit/bruit"
	"sync"

	//"bruit/bruit/clients/kraken/types"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type List interface {
	GetList() List
	AddToEnd(n *Node)
	Print(locker *sync.RWMutex)
	IsEmpty() bool
	GetLast() *Node
}

type WebSocketClient interface {
	// GENERAL METHODS
	InitWebSockets(g *bruit.Settings)
	DeferChanClose(g *bruit.Settings)

	// PUBLIC SOCKET METHODS
	SubscribeToOHLC(g *bruit.Settings, pairs []string, interval int)
	SubscribeToTrades()
	PubDecoder(g *bruit.Settings)
	PubListen(g *bruit.Settings, ohlcMap *OHLCValHolder, tradesWriter *api.WriteAPIImpl) // needs to take an interface instead of OHLCVals
	//SubscribeToTrades()

	// ORDER BOOK SOCKET METHODS
	SubscribeToOrderBook(g *bruit.Settings, pairs []string, depth int)
	//BookDecode(g *ConcurrencySettings)
	BookListen(g *bruit.Settings)

	// PRIVATE SOCKET METHODS

	//PrivListen()
}

type OhlcResponseHolder interface {
	GetChannelID() int
	GetList() *List
	GetInterval() int64
	//Return() *OhlcResponseHolder
}

type OHLCValHolder interface {
	Set(key int, data *OhlcResponseHolder)
	RLock()
	RUnlock()
	Lock()
	Unlock()
	GetVals() map[int]*OhlcResponseHolder
	GetMutex() *sync.RWMutex
}
