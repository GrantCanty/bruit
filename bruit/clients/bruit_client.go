package clients

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
)

type BruitCryptoClient interface {
	//general system commands
	InitClient(s settings.BruitSettings)
	DeferChanClose(s settings.BruitSettings)
	GetHoldingsWithoutStaking() []string
	GetHoldingsWithStaking() []string

	//rest
	GetAssets() (*types.AssetInfoResp, error)
	GetAssetPairs() (*types.AssetPairsResp, error)
	GetOHLC(pair string, interval int) (*types.OHLCResp, error)
	GetAccountBalances() (*types.AccountBalanceResp, error)
	GetPrivateWebSokcetKey() (*types.PrivWSKeyResp, error)

	//ws subscriptions
	SubscribeToTrades(s settings.BruitSettings, pairs []string)
	SubscribeToOHLC(s settings.BruitSettings, pairs []types.Pairs, depth int)
	SubscribeToHoldingsOHLC(s settings.BruitSettings, interval int)
	SubscribeToOrderBook(s settings.BruitSettings, depth int)
	SubscribeToOpenOrders(s settings.BruitSettings, token string)

	//decoders
	PubDecoder(s settings.BruitSettings, OHLCch chan types.OHLCResponse, Tradech chan types.TradeResponse, OHLCSubch chan types.OHLCSuccessResponse)
	BookDecoder(s settings.BruitSettings, Bookch chan types.BookRespV2UpdateJSON, bookDepth int)
	PrivDecoder(s settings.BruitSettings)

	//orders
	CancelAll(s settings.BruitSettings, token string)
	CancelOrder(s settings.BruitSettings, token string, tradeIDs []string)
	AddOrder(s settings.BruitSettings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)

	//Handlers
	HandleOHLCSuccessResponse(resp types.OHLCSuccessResponse)
	HandleOHLCResponse(data types.OHLCResponse, ohlcMap *shared_types.OHLCVals)
}
