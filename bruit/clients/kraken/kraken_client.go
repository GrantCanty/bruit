package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	rest "bruit/bruit/clients/kraken/rest_client"
	"bruit/bruit/clients/kraken/state"
	"bruit/bruit/clients/kraken/types"
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/env"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"errors"
	"log"
)

var ErrConnectionsAlreadInit error = errors.New("connections are already init")

type KrakenClient struct {
	WebSocket web_socket.WebSocketClient
	Rest      rest.RestClient
	State     state.StateManager
}

func (k *KrakenClient) InitClient(s settings.BruitSettings) {
	k.initWebSockets()
	k.initKeys()
	k.initState()
}

func (client *KrakenClient) initWebSockets() {
	if !AreChannelsInit(&client.WebSocket) {
		client.WebSocket.InitChannels()
	}
	client.WebSocket.InitBook()
	client.socketInit()
}

func (k *KrakenClient) initState() error {
	bals, err := k.GetAccountBalances()
	if err != nil {
		return err
	}

	assets, err := k.GetAssets()
	if err != nil {
		return err
	}

	pairs, err := k.GetAssetPairs()
	if err != nil {
		return err
	}

	k.State.Init(*bals, *assets, *pairs)
	return nil
}

// loads the api keys from the .env file
func (k *KrakenClient) initKeys() error {
	env, err := env.Read("CLIENT")
	if err != nil {
		return err
	}
	kraken_data.LoadKeys(env)
	return nil
}

func (client *KrakenClient) socketInit() error {

	// if all sockets are not init, init connections
	if IsPubSocketInit(&client.WebSocket) == nil && IsPrivSocketInit(&client.WebSocket) == nil && IsBookSocketInit(&client.WebSocket) == nil {
		//log.Println("connections are already init")
		return ErrConnectionsAlreadInit
	}
	client.WebSocket.InitSockets()

	// checks to see that sockets are actually init. should switch this to send an error message
	if err := IsPubSocketInit(&client.WebSocket); err != nil { // guard clause checker
		return err
	}
	if err := IsBookSocketInit(&client.WebSocket); err != nil { // guard clause checker
		return err
	}
	if err := IsPrivSocketInit(&client.WebSocket); err != nil { // guard clause checker
		return err
	}

	return nil
}

func (client *KrakenClient) HandleOHLCSuccessResponse(resp types.OHLCSuccessResponse) {
	client.State.Client.AddSubscription(resp.GetMetaData(), types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status})
	//log.Println("subscription list: ", client.State.Client.GetSubscriptions())
}

func (client *KrakenClient) DeferChanClose(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()
	<-s.CtxDone()

	log.Println("Closing channels and connections")

	client.closeChannelsAndConnections()

	log.Println("Closed channels and connections")
}

func (client *KrakenClient) closeChannelsAndConnections() {
	close(client.WebSocket.GetPrivChan())

	if client.WebSocket.GetPubSocket().GetIsConnected() {
		client.WebSocket.GetPubSocket().Close()
	}
	if client.WebSocket.GetBookSocket().GetIsConnected() {
		client.WebSocket.GetBookSocket().Close()
	}
	if client.WebSocket.GetPrivSocket().GetIsConnected() {
		client.WebSocket.GetPrivSocket().Close()
	}
}

func (client *KrakenClient) GetHoldingsWithoutStaking() []string {
	tmp := client.State.Account.GetBalancesWithoutStaking()
	var bals []string
	for i := range tmp {
		bals = append(bals, i)
	}
	return bals
}

func (client *KrakenClient) GetHoldingsWithStaking() []string {
	tmp := client.State.Account.GetBalancesWithStaking()
	var bals []string
	for i := range tmp {
		bals = append(bals, i)
	}
	return bals
}

func (client *KrakenClient) HandleOHLCResponse(data types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	client.State.OnOHLCResponse(data, ohlcMap)
}
