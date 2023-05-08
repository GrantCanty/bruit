package kraken

import (
	"bruit/bruit"
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/clients/kraken/web_socket"
	"bruit/bruit/shared_types"
	"log"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

func (client *KrakenClient) PubListen(g *bruit.Settings, ohlcMap *shared_types.OHLCVals, tradesWriter api.WriteAPI) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	for chanResp := range client.WebSocket.GetPubChan() {
		switch resp := chanResp.(type) {
		case *types.OHLCResponse:
			if g.GlobalSettings.Logging.GetLoggingConsole() {
				log.Println("OHLCResponse")
				log.Printf("new response: %#v\n", resp)
			}
			client.State.OnOHLCResponse(*resp, ohlcMap)
		case *types.TradeResponse:
			if g.GlobalSettings.Logging.GetLoggingConsole() {
				log.Println("TradeResponse")
				log.Printf("new response: %#v %d \n", resp, resp.TradeArray[0].Time.Time.Unix())
			}
			web_socket.OnTradeResponse(*resp, tradesWriter)
			//tradesWriter.WritePoint()
		case *types.ServerConnectionStatusResponse:
			if g.GlobalSettings.Logging.GetLoggingConsole() {
				log.Println("ServerConnectionStatusResponse")
				log.Println(resp)
			}
		case *types.HeartBeat:
			if g.GlobalSettings.Logging.GetLoggingConsole() {
				log.Println("HeartBeat")
				//log.Println(resp.Event)
			}
		case *types.OHLCSuccessResponse:
			if g.GlobalSettings.Logging.GetLoggingConsole() {
				log.Println("OHLCSuccessResponse")
				log.Println(resp)
			}
			client.HandleOHLCSuccessResponse(*resp)
		default:
			log.Println("in default case")
			log.Println(resp)
		}
	}
	<-g.ConcurrencySettings.Ctx.Done()
}

func (client *KrakenClient) BookListen(g *bruit.Settings) {
	g.ConcurrencySettings.Wg.Add(1)
	defer g.ConcurrencySettings.Wg.Done()

	//var book types.BookDecodedResp

	for c := range client.WebSocket.GetBookJSONChan() {
		log.Println("book listen: ", c)
		switch v := c.(type) {
		case *types.HeartBeat:
			log.Println(v)
		case *types.InitialBookResp:
			client.WebSocket.GetBookChan() <- v
		}
	}

	<-g.ConcurrencySettings.Ctx.Done()
	log.Println("closing book listen func")
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
