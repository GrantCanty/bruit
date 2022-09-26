package kraken

import (
	"bruit_new/bruit"
	kraken_data "bruit_new/bruit/clients/kraken/client_data"
	"bruit_new/bruit/clients/kraken/types"
	"bruit_new/bruit/clients/kraken/web_socket"
	"bruit_new/bruit/shared_types"
	"bruit_new/bruit/ws_client"
	"encoding/json"
	"log"
	"strconv"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

// GENERAL METHODS

func (client *KrakenClient) startWebSocketConnection(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	if IsPubSocketInit(client.WebSocket) == nil && IsPrivSocketInit(client.WebSocket) == nil && IsBookSocketInit(client.WebSocket) == nil {
		log.Println("connections are already init")
		return
	}
	client.WebSocket.InitConnections()

	if err := IsPubSocketInit(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}
	if err := IsBookSocketInit(client.WebSocket); err != nil {
		panic(err)
	}
	if err := IsPrivSocketInit(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocketPointer())
	client.WebSocket.GetPubSocketPointer().OnConnected = func(socket ws_client.Socket) {
		log.Println("Connected to public server")
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocketPointer())

	/*ws_client.ReceiveLocker(&client.WebSocket.bookSocket)
	client.WebSocket.bookSocket.OnConnected = func(socket ws_client.Socket) {
		log.Println("Connected to book server")
	}
	ws_client.ReceiveUnlocker(&client.WebSocket.bookSocket)

	ws_client.ReceiveLocker(&client.WebSocket.privSocket)
	client.WebSocket.privSocket.OnConnected = func(socket ws_client.Socket) {
		log.Println("Connected to private server")
	}
	ws_client.ReceiveUnlocker(&client.WebSocket.privSocket)*/

	client.WebSocket.GetPubSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		//decoders.PubJsonDecoder(message, client.Testing)
		client.WebSocket.PubJsonDecoder(message, g.GlobalSettings.Logging)
		log.Println(message)
	}
	/*client.WebSocket.bookSocket.OnTextMessage = func(message string, socket ws_client.Socket) {
		ws_client.BookJsonDecoder(message, client.Testing)
		log.Println(message)
	}
	client.WebSocket.privSocket.OnTextMessage = func(message string, socket ws_client.Socket) {
		ws_client.PrivJsonDecoder(message, client.Testing)
		log.Println(message)
	}*/

	client.WebSocket.GetPubSocketPointer().Connect()
	//client.WebSocket.GetBookSocketPointer().Connect()
	//client.WebSocket.GetPrivSocketPointer().Connect()
	return
}

func (client *KrakenClient) InitWebSockets(g *bruit.Settings) {
	//client.initTesting(testing)
	if !AreChannelsInit(&client.WebSocket) {
		client.WebSocket.InitChannels()
	}
	client.startWebSocketConnection(g)
}

// PUBLIC SOCKET METHODS

func (client *KrakenClient) SubscribeToTrades(g *bruit.Settings, pairs []string) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	if err := PubSocketGuard(client.WebSocket); err != nil {
		panic(err)
	}

	client.WebSocket.SubscribeToTrades(pairs)
}

func (client *KrakenClient) SubscribeToOHLC(g *bruit.Settings, pairs []string, interval int) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	var found bool = false
	for _, i := range kraken_data.OHLCVals {
		if i == interval {
			found = true
			break
		}
	}
	if found == true {
		if err := PubSocketGuard(client.WebSocket); err != nil { // guard clause checker
			panic(err)
		}

		// add func here that makes request to rest OHLC to get past OHLC data. data should then be added to the OHLC map

		client.WebSocket.SubscribeToOHLC(pairs, interval)
	}
}

func (client *KrakenClient) PubDecoder(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	if err := PubSocketGuard(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetPubSocketPointer())
	client.WebSocket.GetPubSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PubJsonDecoder(message, g.GlobalSettings.Logging)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPubSocketPointer())

	<-g.ConcurrencySettings.Ctx.Done()
	return
}

func (client *KrakenClient) PubListen(g *bruit.Settings, ohlcMap shared_types.OHLCValHolder, tradesWriter api.WriteAPI) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	for chanResp := range client.WebSocket.GetPubChan() {
		switch resp := chanResp.(type) {
		case *types.OHLCResponse:
			log.Println("OHLCResponse")
			log.Printf("new response: %#v\n", resp)
			web_socket.OnOHLCResponse(*resp, ohlcMap)
		case *types.TradeResponse:
			log.Println("TradeResponse")
			log.Printf("new response: %#v %d \n", resp, resp.TradeArray[0].Time.Time.Unix())
			web_socket.OnTradeResponse(*resp, tradesWriter)
			//tradesWriter.WritePoint()
		case *types.ServerConnectionStatusResponse:
			log.Println("ServerConnectionStatusResponse")
			log.Println(resp)
		case *types.HeartBeat:
			log.Println("HeartBeat")
			//log.Println(resp.Event)
		case *types.OHLCSuccessResponse:
			log.Println("OHLCSuccessResponse")
			log.Println(resp.ChannelID)
		default:
			log.Println("in default case")
			log.Println(resp)
		}
	}
	<-g.ConcurrencySettings.Ctx.Done()
}

// ORDER BOOK SOCKET METHODS

func (client *KrakenClient) SubscribeToOrderBook(g *bruit.Settings, pairs []string, depth int) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	if err := BookSocketGuard(client.WebSocket); err != nil {
		panic(err)
	}

	sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.BookSubscription{
			Depth: depth,
			Name:  "book",
		},
		Pair: pairs,
	})
	if err != nil {
		log.Println("error marshaling: ", err)
	}
	client.WebSocket.GetBookSocketPointer().SendBinary(sub)
}

func (client *KrakenClient) BookDecoder(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	if err := BookSocketGuard(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}

	ws_client.ReceiveLocker(client.WebSocket.GetBookSocketPointer())
	client.WebSocket.GetBookSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.BookJsonDecoder(message, g.GlobalSettings.Logging)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocketPointer())

	<-g.ConcurrencySettings.Ctx.Done()
	return
}

func (client *KrakenClient) BookListen(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()
	//defer client.WebSocket.GetBookSocketPointer().Close()

	for c := range client.WebSocket.GetBookChan() {
		log.Println(c)
		/*switch v := c.(type) {
		case *types.HeartBeat:
			log.Println(v)
		}*/
	}

	<-g.ConcurrencySettings.Ctx.Done()
	log.Println("closing book listen func")
}

// PRIVATE SOCKET METHODS

func (client *KrakenClient) SubscribeToOpenOrders(g *bruit.Settings, token string) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	PrivSocketGuard(client.WebSocket)

	/*sub, err := json.Marshal(&types.Subscribe{
		Event: "subscribe",
		Subscription: &types.NameAndToken{
			Name:  "openOrders",
			Token: token,
		},
	})

	if err != nil {
		panic(err)
	}

	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)*/
	client.WebSocket.SubscribeToOpenOrders(token)

	//<-g.ConcurrencySettings.Ctx.Done()
}

func (client *KrakenClient) CancelAll(g *bruit.Settings, token string) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	PrivSocketGuard(client.WebSocket)

	sub, _ := json.Marshal(&types.Subscribe{
		Event: "cancelAll",
		Token: token,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)

	<-g.ConcurrencySettings.Ctx.Done()
}

// find a way to ad tradeID
func (client *KrakenClient) CancelOrder(g *bruit.Settings, token string, tradeIDs []string) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	PrivSocketGuard(client.WebSocket)

	sub, _ := json.Marshal(&types.CancelOrder{
		Event: "cancelOrder",
		Token: token,
		Txid:  tradeIDs,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)

	//<-g.ConcurrencySettings.Ctx.Done()
}

func (client *KrakenClient) AddOrder(g *bruit.Settings, token string, otype string, ttype string, pair string, vol string, price string, testing bool) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	PrivSocketGuard(client.WebSocket)

	test := strconv.FormatBool(testing)
	sub, _ := json.Marshal(&types.Order{
		WsToken:   token,
		Event:     "addOrder",
		OrderType: otype,
		TradeType: ttype,
		Pair:      pair,
		Volume:    vol,
		Price:     price,
		Validate:  test,
	})
	client.WebSocket.GetPrivSocketPointer().SendBinary(sub)

	//<-g.ConcurrencySettings.Ctx.Done()
}

func (client *KrakenClient) PrivDecoder(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	PrivSocketGuard(client.WebSocket)

	ws_client.ReceiveLocker(client.WebSocket.GetPrivSocketPointer())
	client.WebSocket.GetPrivSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		client.WebSocket.PrivJsonDecoder(message, g.GlobalSettings.Logging)
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetPrivSocketPointer())

	<-g.ConcurrencySettings.Ctx.Done()
	return
}

func (client *KrakenClient) PrivListen(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	for chanResp := range client.WebSocket.GetPrivChan() {
		switch resp := chanResp.(type) {
		case *types.OpenOrdersResponse:
			log.Println("OpenOrdersResponse")
			log.Println(resp)
		case *types.CancelOrderResponse:
			log.Println("CancelOrderResponse")
			log.Println(resp)
		case *types.ServerConnectionStatusResponse:
			log.Println("ServerConnectionStatusResponse")
			log.Println(resp)
		case *types.HeartBeat:
			log.Println("HeartBeat")
			//log.Println(resp.Event)
		default:
			log.Println("in default case")
			log.Println(resp)
		}
	}
	<-g.ConcurrencySettings.Ctx.Done()
}

// GENERAL METHODS

func (client *KrakenClient) DeferChanClose(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()
	<-g.ConcurrencySettings.Ctx.Done()

	defer close(client.WebSocket.GetPubChan())
	defer close(client.WebSocket.GetBookChan())
	defer close(client.WebSocket.GetPrivChan())

	defer client.WebSocket.GetPubSocketPointer().Close()
	//defer client.WebSocket.GetBookSocketPointer().Close()
	//defer client.WebSocket.GetPrivSocketPointer().Close()

	//ws_client.ReceiveLocker(client.WebSocket.GetBookSocketPointer())
	/*client.WebSocket.GetPubSocketPointer().OnDisconnected = func(err error, socket ws_client.Socket) {
		if err != nil {
			//log.Println("no error: closed pub socke
			log.Println("error: ", err)
		} else {
			log.Println("no error: closed pub socket")
		}
		//client.WebSocket.BookJsonDecoder(message, g.GlobalSettings.Logging)
	}*/
	//ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocketPointer())

	log.Println("Closing channels")
}
