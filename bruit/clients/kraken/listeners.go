package kraken

import (
	"bruit/bruit/clients/kraken/types"
	//web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/settings"
	"log"
)

func (client *KrakenClient) PrivListen(s settings.BruitSettings) {
	s.Add(1)
	defer s.Done()

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
	s.CtxDone()
}
