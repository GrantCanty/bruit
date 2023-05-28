package clients

import "bruit/bruit/settings"

type BruitClient interface {
	InitClient(g settings.Settings)
	DeferChanClose(g settings.Settings)

	SubscribeToTrades(g settings.Settings, pairs []string)
	SubscribeToOHLC(g settings.Settings, pairs []string, depth int)
	SubscribeToOrderBook(g settings.Settings, pairs []string, depth int)
	SubscribeToOpenOrders(g settings.Settings, token string)

	PrivDecoder(g settings.Settings)
	PubDecoder(g settings.Settings)
	BookDecoder(g settings.Settings)

	CancelAll(g settings.Settings, token string)
	CancelOrder(g settings.Settings, token string, tradeIDs []string)
	AddOrder(g settings.Settings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)
}
