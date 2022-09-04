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
	k.InitWebSockets(g)
	k.State.InitState()

}
