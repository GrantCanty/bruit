package kraken

import (
	"bruit/bruit/clients/kraken/types"
	web_socket "bruit/bruit/clients/kraken/web_socket_client"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
	"log"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

func (client *KrakenClient) PubListen(g settings.Settings, ohlcMap *shared_types.OHLCVals, tradesWriter api.WriteAPI) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

	for chanResp := range client.WebSocket.GetPubChan() {
		switch resp := chanResp.(type) {
		case *types.OHLCResponse:
			if g.GetLoggingToConsole() {
				log.Println("OHLCResponse")
				log.Printf("new response: %#v\n", resp)
			}
			client.State.OnOHLCResponse(*resp, ohlcMap)
		case *types.TradeResponse:
			if g.GetLoggingToConsole() {
				log.Println("TradeResponse")
				log.Printf("new response: %#v %d \n", resp, resp.TradeArray[0].Time.Time.Unix())
			}
			web_socket.OnTradeResponse(*resp, tradesWriter)
			//tradesWriter.WritePoint()
		case *types.ServerConnectionStatusResponse:
			if g.GetLoggingToConsole() {
				log.Println("ServerConnectionStatusResponse")
				log.Println(resp)
			}
		case *types.HeartBeat:
			if g.GetLoggingToConsole() {
				log.Println("HeartBeat")
				//log.Println(resp.Event)
			}
		case *types.OHLCSuccessResponse:
			if g.GetLoggingToConsole() {
				log.Println("OHLCSuccessResponse")
				log.Println(resp)
			}
			client.HandleOHLCSuccessResponse(*resp)
		default:
			log.Println("in default case")
			log.Println(resp)
		}
	}
	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
}

func (client *KrakenClient) BookListen(g settings.Settings, book *types.BookDecodedResp) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

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

	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
	log.Println("closing book listen func")
}

func (client *KrakenClient) PrivListen(g settings.Settings) {
	//g.ConcurrencySettings.Wg.Add(1)
	g.Add(1)
	//defer g.ConcurrencySettings.Wg.Done()
	defer g.Done()

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
	//<-g.ConcurrencySettings.Ctx.Done()
	//<-g.CtxDone()
	g.CtxDone()
}
