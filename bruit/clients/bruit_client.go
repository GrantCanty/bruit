package clients

import (
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type BruitClient interface {
	InitClient(g settings.Settings)
	DeferChanClose(g settings.Settings)

	SubscribeToTrades(g settings.Settings, pairs []string)
	SubscribeToOHLC(g settings.Settings, pairs []string, depth int)
	SubscribeToOrderBook(g settings.Settings, pairs []string, depth int)
	SubscribeToOpenOrders(g settings.Settings, token string)

	PubDecoder(g settings.Settings)
	BookDecoder(g settings.Settings)
	PrivDecoder(g settings.Settings)

	PubListen(g settings.Settings, ohlcMap *shared_types.OHLCVals, tradesWriter api.WriteAPI)

	CancelAll(g settings.Settings, token string)
	CancelOrder(g settings.Settings, token string, tradeIDs []string)
	AddOrder(g settings.Settings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)
}
