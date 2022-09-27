package web_socket

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"

	//"container/list"
	"strconv"
	"time"

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

func OnOHLCResponse(data types.OHLCResponse, ohlcMap shared_types.OHLCValHolder) {
	/**
	*  Add:
	*  OHLCResponseHandler func to add responses to a LL. should delete the head if length is too long (ex: 10,000)
	*  CalcTechnicals func to recalculate the values of technical indicators
	*  Eval func to evaluate if buy/sell condition is met
	*  PlaceOrder func depending on Eval func
	**/
	HandleOHLCResponse(data, ohlcMap)
	//oo := ohlcMap.Vals
	//gg := (*(*ohlcMap).GetVals()[data.ChannelID]).GetList()
	//fmt.Println(gg)

	ohlcMap.GetData()[data.ChannelID].GetList().Print(ohlcMap.GetMutex())
}

func HandleOHLCResponse(v types.OHLCResponse, ohlcMap shared_types.OHLCValHolder) {
	ohlcMap.RLock()
	if subID, found := ohlcMap.GetData()[v.ChannelID]; found { // if channelID already exists in the map, then...

		var tmpStartTime time.Time = v.OHLCArray.EndTime.Add(-time.Minute * time.Duration(subID.GetInterval()))
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		v.OHLCArray.SetStartTime(tmpStartTime)
		ohlcMap.Unlock()

		ohlcMap.RLock()
		if subID.GetList().IsEmpty() { // if no responses
			node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			subID.GetList().AddToEnd(&node)
			ohlcMap.Unlock()
		} else if subID.GetList().GetLast().Data.GetEndTime().Time.Equal(v.OHLCArray.EndTime.Time) { // if updating last candle
			tmp := v.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			//editCandle(subID.GetList().GetLast().Data.(*types.WSCandles), tmp) // pay attention to this since i believe it should be passed as ref, not as val
			subID.GetList().EditCandle(subID.GetList().GetLast().Data, &tmp)
			ohlcMap.Unlock()
		} else if subID.GetList().GetLast().Data.GetEndTime().Time.Equal(v.OHLCArray.StartTime.Time) { // adding candle to next index in array
			node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			subID.GetList().AddToEnd(&node)
			ohlcMap.Unlock()
		} else if subID.GetList().GetLast().Data.GetEndTime().Time.Before(v.OHLCArray.StartTime.Time) { // if adding multiple candles
			tmp := v.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			//addCandle(*subID.GetList(), tmp, subID.GetInterval())
			subID.GetList().AddCandle(&tmp, &types.WSCandles{}, subID.GetInterval())
			ohlcMap.Unlock()
		}
	} else { // if the channel id cannot be found in the map
		interval, _ := strconv.ParseInt(v.ChannelName[len(v.ChannelName)-1:], 10, 64)
		v.OHLCArray.StartTime.Time = v.OHLCArray.EndTime.Add(-time.Minute * time.Duration(interval))
		node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
		tmp := shared_types.OhlcResponseHolder(&types.KrakenOHLCResponseHolder{ChannelID: v.ChannelID, ChannelName: v.ChannelName, Pair: v.Pair, Interval: interval, List: shared_types.List{Head: &node, Last: &node, Length: 1}})
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		ohlcMap.Set(tmp.GetChannelID(), tmp)
		ohlcMap.Unlock()
	}
}

/*func editCandle(oldCandle *types.WSCandles, newCandle types.WSCandles) {
	if oldCandle.GetHigh().LessThan(newCandle.High) {
		oldCandle.SetHigh(newCandle.High) // = newCandle.High
	}
	if oldCandle.GetLow().GreaterThan(newCandle.Low) {
		oldCandle.SetLow(newCandle.Low) // = newCandle.Low
	}
	oldCandle.SetClose(newCandle.Close)   // = newCandle.Close
	oldCandle.SetVWAP(newCandle.VWAP)     // = newCandle.VWAP
	oldCandle.SetVolume(newCandle.Volume) // = newCandle.Volume
	oldCandle.SetCount(newCandle.Count)   // = newCandle.Count
}*/

/*func addCandle(list shared_types.List, newCandle types.WSCandles, interval int64) { // old candle should switch to list
	since := time.Since(list.GetLast().Data.GetEndTime().Time).Minutes()
	if since < time.Duration(interval).Minutes() { // if the time since the close of the last candle is less than the time of the interval, the candle you received will just be added to the end
		newCandle.StartTime = list.GetLast().Data.GetEndTime()
		node := shared_types.Node{Data: &newCandle, Next: nil}
		list.AddToEnd(&node)
	} else {
		newNodeCount := int64(int64(since) / interval)
		zero := decimal.New(0, 0)
		for i := int64(0); i < newNodeCount; i++ {
			last := list.GetLast()
			nodeData := types.WSCandles{
				StartTime: last.Data.GetStartTime(),
				EndTime:   last.Data.GetEndTime(),
				Open:      last.Data.GetClose(),
				High:      last.Data.GetClose(),
				Low:       last.Data.GetClose(),
				Close:     last.Data.GetClose(),
				VWAP:      zero,
				Volume:    zero,
				Count:     0,
			}
			nodeData.StartTime.Time = nodeData.StartTime.Time.Add(time.Minute * time.Duration(interval))
			nodeData.EndTime.Time = nodeData.EndTime.Time.Add(time.Minute * time.Duration(interval))
			node := shared_types.Node{Data: &nodeData, Next: nil}
			list.AddToEnd(&node)
		}
		node := shared_types.Node{Data: &newCandle, Next: nil}
		list.AddToEnd(&node)
	}
}*/
