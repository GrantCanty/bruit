package state

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
	"strconv"
	"time"
)

func (sm *StateManager) Init(bals types.AccountBalanceResp, assets types.AssetInfoResp, pairs types.AssetPairsResp) {
	sm.Client.initClient()

	sm.Account.initAccount(bals)
	sm.Client.initAssets(assets)
	sm.Client.initPairs(pairs)
}

func (sm *StateManager) GetSubscriptions() map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData {
	return sm.Client.subscriptions
}

func (sm *StateManager) GetInterval(metaData shared_types.SubscriptionMetaData) int {
	//log.Println("interval: ", sm.Client.subscriptions[metaData])
	return sm.Client.subscriptions[metaData].GetData().(types.KrakenOHLCSubscriptionData).Interval
}

func (sm *StateManager) GetChannelID(metaData shared_types.SubscriptionMetaData) int {
	return metaData.GetChannelID()
}

func (sm *StateManager) OnOHLCResponse(data types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	sm.handleOHLCResponse(data, ohlcMap)
	ohlcMap.GetData()[data.GetMetaData()].GetList().Print(ohlcMap.GetMutex())
}

func (sm *StateManager) handleOHLCResponse(resp types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	ohlcMap.Lock()
	defer ohlcMap.Unlock()

	if list, found := ohlcMap.GetData()[resp.GetMetaData()]; found { // if channelID already exists in the map, then...
		var tmpStartTime time.Time = resp.OHLCArray.EndTime.Add(-time.Minute * time.Duration(sm.GetInterval(resp.GetMetaData())))
		resp.OHLCArray.SetStartTime(tmpStartTime)

		if list.GetList().IsEmpty() { // if no responses
			node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
			list.GetList().AddToEnd(&node)
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Equal(resp.OHLCArray.EndTime.Time) { // if updating last candle
			tmp := resp.OHLCArray
			list.GetList().EditCandle(list.GetList().GetLast().Data, &tmp)
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Equal(resp.OHLCArray.StartTime.Time) { // adding candle to next index in array
			node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
			list.GetList().AddToEnd(&node)
		} else if list.GetList().GetLast().Data.GetEndTime().Time.Before(resp.OHLCArray.StartTime.Time) { // if adding multiple candles
			tmp := resp.OHLCArray
			list.GetList().AddCandle(&tmp, &types.WSCandles{}, sm.GetInterval(resp.GetMetaData()))
		}
	} else { // if the channel id cannot be found in the map
		interval, _ := strconv.ParseInt(resp.ChannelName[len(resp.ChannelName)-1:], 10, 64)
		resp.OHLCArray.StartTime.Time = resp.OHLCArray.EndTime.Add(-time.Minute * time.Duration(interval))
		node := shared_types.Node{Data: &resp.OHLCArray, Next: nil}
		tmp := &shared_types.List{Head: &node, Last: &node, Length: 1}
		ohlcMap.Set(resp.GetMetaData(), tmp)
	}
}
