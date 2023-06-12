package clients

import (
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type BruitClient interface {
	InitClient(g settings.BruitSettings)
	DeferChanClose(g settings.BruitSettings)

	SubscribeToTrades(g settings.BruitSettings, pairs []string)
	SubscribeToOHLC(g settings.BruitSettings, pairs []string, depth int)
	SubscribeToOrderBook(g settings.BruitSettings, pairs []string, depth int)
	SubscribeToOpenOrders(g settings.BruitSettings, token string)

	PubDecoder(g settings.BruitSettings)
	BookDecoder(g settings.BruitSettings)
	PrivDecoder(g settings.BruitSettings)

	PubListen(g settings.BruitSettings, ohlcMap *shared_types.OHLCVals, tradesWriter api.WriteAPI)

	CancelAll(g settings.BruitSettings, token string)
	CancelOrder(g settings.BruitSettings, token string, tradeIDs []string)
	AddOrder(g settings.BruitSettings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)
}
