package shared_types

import (
	"bruit_new/bruit"
	"sync"

	//"bruit/bruit/clients/kraken/types"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

/*type List interface {
	GetList() List
	AddToEnd(n *Node)
	Print(locker *sync.RWMutex)
	IsEmpty() bool
	GetLast() *Node
}*/

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
