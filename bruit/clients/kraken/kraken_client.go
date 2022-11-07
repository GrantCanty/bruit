package kraken

import (
	"bruit/bruit"
	"bruit/bruit/clients/kraken/rest"
	"bruit/bruit/clients/kraken/state"
	"bruit/bruit/clients/kraken/web_socket"
)

type KrakenClient struct {
	WebSocket web_socket.WebSocketClient
	Rest      rest.RestClient
	State     state.StateManager
	//Engine
}

func (k *KrakenClient) InitClient(g *bruit.Settings) {
	k.initWebSockets(g)
	k.initState()

}

func (k *KrakenClient) initState() {
	bals, err := k.GetAccountBalances()
	if err != nil {
		panic(err)
	}
	k.State.Init(*bals)
}
