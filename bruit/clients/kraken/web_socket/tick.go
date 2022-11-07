package web_socket

import (
	"bruit/bruit/clients/kraken/types"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func OnTradeResponse(data types.TradeResponse, tradesWriter api.WriteAPI) {
	for _, trade := range data.TradeArray {
		price, _ := trade.Price.Float64()
		vol, _ := trade.Price.Float64()
		//influxdb2
		point := influxdb2.NewPoint("trade", map[string]string{"object": "trade", "pair": data.Pair}, map[string]interface{}{"price": price, "volume": vol, "side": trade.Side, "orderType": trade.OrderType, "misc": trade.Misc}, trade.Time.Time)
		tradesWriter.WritePoint(point)
	}
}

/*func (ws *WebSocketClient) OnOHLCResponse(data types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	ws.handleOHLCResponse(data, ohlcMap)
	//oo := ohlcMap.Vals
	//gg := (*(*ohlcMap).GetVals()[data.ChannelID]).GetList()
	//fmt.Println(gg)

	ohlcMap.GetData()[data.GetMetaData()].GetList().Print(ohlcMap.GetMutex())
}

func (ws *WebSocketClient) handleOHLCResponse(resp types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	ohlcMap.RLock()
	if list, found := ohlcMap.GetData()[resp.GetMetaData()]; found { // if channelID already exists in the map, then...

		var tmpStartTime time.Time = resp.OHLCArray.EndTime.Add(-time.Minute * time.Duration(ws.GetInterval(resp.GetMetaData())))
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		resp.OHLCArray.SetStartTime(tmpStartTime)
		ohlcMap.Unlock()

		ohlcMap.RLock()
		if list.GetList().IsEmpty() { // if no responses
			node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			list.GetList().AddToEnd(&node)
			ohlcMap.Unlock()
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Equal(resp.OHLCArray.EndTime.Time) { // if updating last candle
			tmp := resp.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			list.GetList().EditCandle(list.GetList().GetLast().Data, &tmp)
			ohlcMap.Unlock()
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Equal(resp.OHLCArray.StartTime.Time) { // adding candle to next index in array
			node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			list.GetList().AddToEnd(&node)
			ohlcMap.Unlock()
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Before(resp.OHLCArray.StartTime.Time) { // if adding multiple candles
			tmp := resp.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			list.GetList().AddCandle(&tmp, &types.WSCandles{}, ws.GetInterval(resp.GetMetaData()))
			ohlcMap.Unlock()
		}
	} else { // if the channel id cannot be found in the map
		interval, _ := strconv.ParseInt(resp.ChannelName[len(resp.ChannelName)-1:], 10, 64)
		resp.OHLCArray.StartTime.Time = resp.OHLCArray.EndTime.Add(-time.Minute * time.Duration(interval))
		node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
		tmp := &shared_types.List{Head: &node, Last: &node, Length: 1}
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		ohlcMap.Set(resp.GetMetaData(), tmp)
		ohlcMap.Unlock()
	}
}*/
