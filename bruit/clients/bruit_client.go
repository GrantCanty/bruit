package clients

import "bruit/bruit"

type BruitClient interface {
	InitClient(g *bruit.Settings)
	DeferChanClose(g *bruit.Settings)

	SubscribeToTrades(g *bruit.Settings, pairs []string)
	SubscribeToOHLC(g *bruit.Settings, pairs []string, depth int)
	SubscribeToOrderBook(g *bruit.Settings, pairs []string, depth int)
	SubscribeToOpenOrders(g *bruit.Settings, token string)

	PrivDecoder(g *bruit.Settings)
	PubDecoder(g *bruit.Settings)
	BookDecoder(g *bruit.Settings)

	CancelAll(g *bruit.Settings, token string)
	CancelOrder(g *bruit.Settings, token string, tradeIDs []string)
	AddOrder(g *bruit.Settings, token string, otype string, ttype string, pair string, vol string, price string, testing bool)
}
