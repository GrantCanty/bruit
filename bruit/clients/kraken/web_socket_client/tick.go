package web_socket

import (
	"bruit/bruit/clients/kraken/types"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func OnTradeResponse(data types.TradeResponse, tradesWriter api.WriteAPI) {
	for _, trade := range data.TradeArray {
		price, _ := trade.Price.Float64()
		vol, _ := trade.Volume.Float64()
		//influxdb2
		point := influxdb2.NewPoint("trade", map[string]string{"object": "trade", "pair": data.Pair}, map[string]interface{}{"price": price, "volume": vol, "side": trade.Side, "orderType": trade.OrderType, "misc": trade.Misc}, trade.Time.Time)
		tradesWriter.WritePoint(point)
	}
}
