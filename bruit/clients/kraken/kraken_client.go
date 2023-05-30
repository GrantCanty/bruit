package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	rest "bruit/bruit/clients/kraken/rest_client"
	"bruit/bruit/clients/kraken/state"
	"bruit/bruit/clients/kraken/types"
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/env"
	"bruit/bruit/settings"
	"bruit/bruit/ws_client"
	"log"
)

type KrakenClient struct {
	WebSocket web_socket.WebSocketClient
	Rest      rest.RestClient
	State     state.StateManager
}

func (k *KrakenClient) InitClient(g settings.Settings) {
	k.initWebSockets(g)
	k.initKeys()
	k.initState()

}

func (client *KrakenClient) initWebSockets(g settings.Settings) {
	if !AreChannelsInit(&client.WebSocket) {
		client.WebSocket.InitChannels()
	}
	client.startWebSocketConnection(g)
}

func (k *KrakenClient) initState() {
	bals, err := k.GetAccountBalances()
	if err != nil {
		panic(err)
	}
	k.State.Init(*bals)
}

//loads the api keys from the .env file
func (k *KrakenClient) initKeys() {
	env, err := env.Read()
	if err != nil {
		panic(err)
	}
	kraken_data.LoadKeys(env)
}

func (client *KrakenClient) startWebSocketConnection(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	//g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	//defer g.Done()

	// if all sockets are not init, init connections
	if IsPubSocketInit(client.WebSocket) == nil && IsPrivSocketInit(client.WebSocket) == nil && IsBookSocketInit(client.WebSocket) == nil {
		log.Println("connections are already init")
		return
	}
	client.WebSocket.InitConnections()

	// checks to see that sockets are actually init. should switch this to send an error message
	if err := IsPubSocketInit(client.WebSocket); err != nil { // guard clause checker
		panic(err)
	}
	if err := IsBookSocketInit(client.WebSocket); err != nil { // guard clause checker
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

	/*ws_client.ReceiveLocker(client.WebSocket.GetBookSocketPointer())
	client.WebSocket.GetBookSocketPointer().OnConnected = func(socket ws_client.Socket) {
		log.Println("Connected to book server")
	}
	ws_client.ReceiveUnlocker(client.WebSocket.GetBookSocketPointer())*/

	/*ws_client.ReceiveLocker(&client.WebSocket.privSocket)
	client.WebSocket.privSocket.OnConnected = func(socket ws_client.Socket) {
		log.Println("Connected to private server")
	}
	ws_client.ReceiveUnlocker(&client.WebSocket.privSocket)*/

	/*client.WebSocket.GetPubSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		//decoders.PubJsonDecoder(message, client.Testing)
		client.WebSocket.PubJsonDecoder(message, g.GlobalSettings.Logging)
		log.Println(message)
	}*/

	/*client.WebSocket.GetBookSocketPointer().OnTextMessage = func(message string, socket ws_client.Socket) {
		//decoders.BookJsonDecoder(message, client.Testing)
		client.WebSocket.BookJsonDecoder(message, g.GlobalSettings.Logging)
		log.Println(message)
	}*/

	/*client.WebSocket.privSocket.OnTextMessage = func(message string, socket ws_client.Socket) {
		ws_client.PrivJsonDecoder(message, client.Testing)
		log.Println(message)
	}*/

	client.WebSocket.GetPubSocketPointer().Connect()
	//client.WebSocket.GetBookSocketPointer().Connect()
	//client.WebSocket.GetPrivSocketPointer().Connect()
	//<-g.CtxDone()
	return
}

func (client *KrakenClient) HandleOHLCSuccessResponse(resp types.OHLCSuccessResponse) {
	client.State.Client.AddSubscription(resp.GetMetaData(), types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status})
	//log.Println("subscription list: ", client.State.Client.GetSubscriptions())
}

func (client *KrakenClient) DeferChanClose(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()
	//<-g.ConcurrencySettings.Ctx.Done()
	g.CtxDone()

	log.Println("Closing channels")

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

	log.Println("Closed channels")
}
