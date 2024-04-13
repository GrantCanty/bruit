package kraken

import (
	"bruit/bruit/clients/kraken/types"
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/settings"
	"log"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

func (client *KrakenClient) PubListen(s settings.BruitSettings, ch chan types.OHLCResponse, tradesWriter api.WriteAPI) {
	s.Add(1)
	defer s.Done()

	for chanResp := range client.WebSocket.GetPubChan() {
		switch resp := chanResp.(type) {
		case *types.OHLCResponse:
			if s.GetLoggingToConsole() {
				//log.Println("OHLCResponse")
				//log.Printf("new response: %#v\n", resp)
			}
			//client.State.OnOHLCResponse(*resp, ohlcMap)
			ch <- *resp
		case *types.TradeResponse:
			if s.GetLoggingToConsole() {
				log.Println("TradeResponse")
				log.Printf("new response: %#v %d \n", resp, resp.TradeArray[0].Time.Time.Unix())
			}
			web_socket.OnTradeResponse(*resp, tradesWriter)
			//tradesWriter.WritePoint()
		case *types.ServerConnectionStatusResponse:
			if s.GetLoggingToConsole() {
				log.Println("ServerConnectionStatusResponse")
				log.Println(resp)
			}
		case *types.HeartBeat:
			if s.GetLoggingToConsole() {
				log.Println("HeartBeat")
				//log.Println(resp.Event)
			}
		case *types.OHLCSuccessResponse:
			if s.GetLoggingToConsole() {
				log.Println("OHLCSuccessResponse")
				log.Println(resp)
			}
			client.HandleOHLCSuccessResponse(*resp)
		default:
			log.Println("in default case")
			log.Println(resp)
		}
	}
	s.CtxDone()
}

func (client *KrakenClient) BookListen(s settings.BruitSettings, book *types.BookDecodedResp) {
	s.Add(1)
	defer s.Done()

	//var book types.BookDecodedResp

	for c := range client.WebSocket.GetBookJSONChan() {
		log.Println("book listen: ", c)
		switch v := c.(type) {
		case *types.HeartBeat:
			log.Println(v)
		case *types.BookDecodedResp:
			book = v
		}
	}

	s.CtxDone()
	log.Println("closing book listen func")
}

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
