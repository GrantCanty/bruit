package web_socket

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/shopspring/decimal"
)

func OnOHLCResponse(data *types.OHLCResponse, ohlcMap *types.OHLCVals) {
	/**
	*  Add:
	*  OHLCResponseHandler func to add responses to a LL. should delete the head if length is too long (ex: 10,000)
	*  CalcTechnicals func to recalculate the values of technical indicators
	*  Eval func to evaluate if buy/sell condition is met
	*  PlaceOrder func depending on Eval func
	**/
	HandleOHLCResponse(data, ohlcMap)
	ohlcMap.Vals[data.ChannelID].List.Print(&ohlcMap.Mutex)
}

func OnTradeResponse(data types.TradeResponse, tradesWriter api.WriteAPI) {
	for _, trade := range data.TradeArray {
		price, _ := trade.Price.Float64()
		vol, _ := trade.Price.Float64()
		point := influxdb2.NewPoint("trade", map[string]string{"object": "trade", "pair": data.Pair}, map[string]interface{}{"price": price, "volume": vol, "side": trade.Side, "orderType": trade.OrderType, "misc": trade.Misc}, trade.Time.Time)
		tradesWriter.WritePoint(point)
	}
}

func HandleOHLCResponse(v *types.OHLCResponse, ohlcMap *types.OHLCVals) {
	ohlcMap.RLock()
	if subID, found := ohlcMap.Vals[v.ChannelID]; found { // if channelID already exists in the map, then...

		var tmpStartTime time.Time = v.OHLCArray.EndTime.Add(-time.Minute * time.Duration(subID.Interval))
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		v.OHLCArray.NewStartTime(tmpStartTime)
		ohlcMap.Unlock()

		ohlcMap.RLock()
		if subID.List.Empty() { // if no responses
			node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			subID.List.AddToEnd(&node)
			ohlcMap.Unlock()
		} else if subID.List.Last.Data.GetEndTime().Time.Equal(v.OHLCArray.EndTime.Time) { // if updating last candle
			tmp := v.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			editCandle(subID.List.Last.Data.(*types.WSCandles), tmp) // pay attention to this since i believe it should be passed as ref, not as val
			ohlcMap.Unlock()
		} else if subID.List.Last.Data.GetEndTime().Time.Equal(v.OHLCArray.StartTime.Time) { // adding candle to next index in array
			//data := v.OHLCArray
			node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			subID.List.AddToEnd(&node)
			ohlcMap.Unlock()
		} else if subID.List.Last.Data.GetEndTime().Time.Before(v.OHLCArray.StartTime.Time) { // if adding multiple candles
			tmp := v.OHLCArray
			ohlcMap.RUnlock()

			ohlcMap.Lock()
			addCandle(&subID.List, tmp, subID.Interval)
			ohlcMap.Unlock()
		}
	} else { // if the channel id cannot be found in the map
		interval, _ := strconv.ParseInt(v.ChannelName[len(v.ChannelName)-1:], 10, 64)
		v.OHLCArray.StartTime.Time = v.OHLCArray.EndTime.Add(-time.Minute * time.Duration(interval))
		node := shared_types.Node{Data: &v.OHLCArray, Next: nil}
		tmp := types.OHLCResponseHolder{ChannelID: v.ChannelID, ChannelName: v.ChannelName, Pair: v.Pair, Interval: interval, List: shared_types.List{Head: &node, Last: &node, Length: 1}}
		ohlcMap.RUnlock()

		ohlcMap.Lock()
		ohlcMap.Set(tmp.ChannelID, &tmp)
		ohlcMap.Unlock()
	}
}

func editCandle(oldCandle *types.WSCandles, newCandle types.WSCandles) {
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
}

func addCandle(list *shared_types.List, newCandle types.WSCandles, interval int64) { // old candle should switch to list
	since := time.Since(list.Last.Data.GetEndTime().Time).Minutes()
	if since < time.Duration(interval).Minutes() { // if the time since the close of the last candle is less than the time of the interval, the candle you received will just be added to the end
		newCandle.StartTime = list.Last.Data.GetEndTime()
		node := shared_types.Node{Data: &newCandle, Next: nil}
		list.AddToEnd(&node)
	} else {
		newNodeCount := int64(int64(since) / interval)
		zero := decimal.New(0, 0)
		for i := int64(0); i < newNodeCount; i++ {
			nodeData := types.WSCandles{
				StartTime: list.Last.Data.GetStartTime(),
				EndTime:   list.Last.Data.GetEndTime(),
				Open:      list.Last.Data.GetClose(),
				High:      list.Last.Data.GetClose(),
				Low:       list.Last.Data.GetClose(),
				Close:     list.Last.Data.GetClose(),
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
}
