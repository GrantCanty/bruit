package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	rest "bruit/bruit/clients/kraken/rest_client"
	"bruit/bruit/clients/kraken/state"
	"bruit/bruit/clients/kraken/types"
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/env"
	"bruit/bruit/settings"
	"log"
)

type KrakenClient struct {
	WebSocket web_socket.WebSocketClient
	Rest      rest.RestClient
	State     state.StateManager
}

func (k *KrakenClient) InitClient(g settings.BruitSettings) {
	k.initWebSockets()
	k.initKeys()
	k.initState()

}

func (client *KrakenClient) initWebSockets() {
	if !AreChannelsInit(&client.WebSocket) {
		client.WebSocket.InitChannels()
	}
	client.socketInit()
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
	env, err := env.Read("CLIENT")
	if err != nil {
		panic(err)
	}
	kraken_data.LoadKeys(env)
}

func (client *KrakenClient) socketInit() {

	// if all sockets are not init, init connections
	if IsPubSocketInit(client.WebSocket) == nil && IsPrivSocketInit(client.WebSocket) == nil && IsBookSocketInit(client.WebSocket) == nil {
		log.Println("connections are already init")
		return
	}
	client.WebSocket.InitSockets()

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
	return
}

func (client *KrakenClient) HandleOHLCSuccessResponse(resp types.OHLCSuccessResponse) {
	client.State.Client.AddSubscription(resp.GetMetaData(), types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status})
	//log.Println("subscription list: ", client.State.Client.GetSubscriptions())
}

func (client *KrakenClient) DeferChanClose(g settings.BruitSettings) {
	g.Add(1)
	defer g.Done()
	<-g.CtxDone()

	log.Println("Closing channels and connections")

	client.closeChannelsAndConnections()

	log.Println("Closed channels and connections")
}

func (client *KrakenClient) closeChannelsAndConnections() {
	close(client.WebSocket.GetPubChan())
	close(client.WebSocket.GetBookChan())
	close(client.WebSocket.GetPrivChan())

	if client.WebSocket.GetPubSocket().IsConnected {
		client.WebSocket.GetPubSocketPointer().Close()
	}
	if client.WebSocket.GetBookSocket().IsConnected {
		client.WebSocket.GetBookSocketPointer().Close()
	}
	if client.WebSocket.GetPrivSocket().IsConnected {
		client.WebSocket.GetPrivSocketPointer().Close()
	}
}
