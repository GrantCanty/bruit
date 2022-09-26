package kraken

import (
	"bruit_new/bruit"
	"bruit_new/bruit/clients/kraken/rest"
	"bruit_new/bruit/clients/kraken/state"
	"bruit_new/bruit/clients/kraken/web_socket"
)

type KrakenClient struct {
	WebSocket web_socket.WebSocketClient
	Rest      rest.RestClient
	State     state.StateManager
	//Engine
}

func (k *KrakenClient) InitClient(g *bruit.Settings) {
	k.InitWebSockets(g)
	k.State.InitState()

}
