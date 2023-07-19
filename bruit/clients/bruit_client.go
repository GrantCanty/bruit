package clients

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type BruitClient interface {
	//general system commands
	InitClient(g settings.BruitSettings)
	DeferChanClose(g settings.BruitSettings)
	GetHoldingsWithoutStaking() []string
	GetHoldingsWithStaking() []string

	//rest
	GetAssets() (*types.AssetInfoResp, error)
	GetAssetPairs() (*types.AssetPairsResp, error)
	GetOHLC(pair string, interval int) (*types.OHLCResp, error)
	GetAccountBalances() (*types.AccountBalanceResp, error)
	GetPrivateWebSokcetKey() (*types.PrivWSKeyResp, error)

	//ws subscriptions
	SubscribeToTrades(g settings.BruitSettings, pairs []string)
	SubscribeToOHLC(g settings.BruitSettings, pairs []types.Pairs, depth int)
	SubscribeToHoldingsOHLC(g settings.BruitSettings, interval int)
	SubscribeToOrderBook(g settings.BruitSettings, pairs []string, depth int)
	SubscribeToOpenOrders(g settings.BruitSettings, token string)

	//decoders
	PubDecoder(g settings.BruitSettings)
	BookDecoder(g settings.BruitSettings)
	PrivDecoder(g settings.BruitSettings)

	//listeners
	PubListen(g settings.BruitSettings, ohlcMap *shared_types.OHLCVals, tradesWriter api.WriteAPI)

	//orders
	CancelAll(g settings.BruitSettings, token string)
	CancelOrder(g settings.BruitSettings, token string, tradeIDs []string)
	AddOrder(g settings.BruitSettings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)
}
