package main

import (
	"bruit/bruit"
	"bruit/bruit/clients/kraken"
	"bruit/bruit/influx"
	"bruit/bruit/shared_types"
)

func main() {
	g := bruit.Settings{}
	g.Init()

	db := influx.DB{}
	db.Init()
	//balancesWriter := db.WriteAPIBlocking("Vert", "Balances")

	k := &kraken.KrakenClient{}
	k.InitClient(&g)
	/*balances, err := k.GetAccountBalance()
	if err != nil {
		panic(err)
	}
	bals := make(map[string]interface{})
	for pair, vol := range *balances {
		log.Println(pair, ": ", (*balances)[pair])
		if vol.GreaterThan(decimal.NewFromFloat(0.000001)) {
			val, _ := vol.Float64()
			bals[pair] = float64(val)
		}
	}
	point := influxdb2.NewPoint("balances", map[string]string{"object": "newBalance"}, bals, time.Now())
	log.Print(*balances)
	log.Println(bals)
	log.Println(point)
	balancesWriter.WritePoint(ctx, point)*/

	/*res, err := k.GetAssets()
	if err != nil {
		panic(err)
	}
	log.Println(res)

	resp, err := k.GetAssetPairs()
	if err != nil {
		panic(err)
	}
	log.Println(resp)

	respp, err := k.GetOHLC("XXBTZUSD", 5)
	if err != nil {
		panic(err)
	}
	log.Println(respp)

	resp, err := k.GetPrivateWebSokcetKey()
	if err != nil {
		panic(err)
	}*/
	//log.Println(resp.Token)

	//k.InitWebSockets(&g)

	go k.PubDecoder(&g)
	//go k.BookDecoder(&g)
	//go k.PrivDecoder(&g)
	//go k.PrivListen(&g)

	ohlcMap := shared_types.OHLCVals{}
	go k.PubListen(&g, &ohlcMap, db.GetTradeWriter())
	//go k.BookListen(&g)

	//k.SubscribeToTrades(&g, []string{"BTC/USD", "ETH/USD"})
	k.SubscribeToOHLC(&g, []string{"EOS/USD", "BTC/USD"}, 1)
	//k.SubscribeToOrderBookk(g, []string{"BTC/USD"}, 10)
	//go k.SubscribeToOpenOrders(&g, resp.Token)*/
	//k.SubscribeToOrderBook(&g, []string{"EOS/USD"}, 10)

	/*bots := []shared_types.WebSocketClient{k}
	for _, bot := range bots {
		fmt.Println(bot)
	}*/

	//go k.BookDecodee(g)
	//go k.BookListenn(g)
	//go k.PrivListen(&g)

	//k.WebSocket.PrivChan
	go k.DeferChanClose(&g)
	g.Wait()

}
