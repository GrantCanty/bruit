package state

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
	"log"
	"strconv"
	"time"
)

func (sm *StateManager) Init(bals types.AccountBalanceResp, assets types.AssetInfoResp, pairs types.AssetPairsResp) {
	sm.Client.initClient()

	sm.Account.initAccount(bals)
	sm.Client.initAssets(assets)
	sm.Client.initPairs(pairs)
}

func (sm StateManager) GetSubscriptions() map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData {
	return sm.Client.subscriptions
}

func (sm StateManager) GetInterval(metaData shared_types.SubscriptionMetaData) int {
	//log.Println("interval: ", sm.Client.subscriptions[metaData])
	return sm.Client.subscriptions[metaData].GetData().(types.KrakenOHLCSubscriptionData).Interval
}

func (sm StateManager) GetChannelID(metaData shared_types.SubscriptionMetaData) int {
	return metaData.GetChannelID()
}

func (sm *StateManager) HandleOHLCSuccessResponse(resp *types.OHLCSuccessResponse) {
	if _, found := sm.Client.subscriptions[resp.GetMetaData()]; found {
		sm.Client.subscriptions[resp.GetMetaData()] = types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status}
	} else {
		if sm.Client.subscriptions == nil {
			sm.Client.subscriptions = make(map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData)
		}
		sm.Client.subscriptions[resp.GetMetaData()] = types.KrakenOHLCSubscriptionData{Interval: resp.Subscription.Interval, Status: resp.Status}
	}
	log.Println("subscription list: ", sm.Client.subscriptions)
}

func (sm *StateManager) OnOHLCResponse(data types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	/**
	*  Add:
	*  OHLCResponseHandler func to add responses to a LL. should delete the head if length is too long (ex: 10,000)
	*  CalcTechnicals func to recalculate the values of technical indicators
	*  Eval func to evaluate if buy/sell condition is met
	*  PlaceOrder func depending on Eval func
	**/
	sm.handleOHLCResponse(data, ohlcMap)
	//oo := ohlcMap.Vals
	//gg := (*(*ohlcMap).GetVals()[data.ChannelID]).GetList()
	//fmt.Println(gg)

	ohlcMap.GetData()[data.GetMetaData()].GetList().Print(ohlcMap.GetMutex())
}

func (sm *StateManager) handleOHLCResponse(resp types.OHLCResponse, ohlcMap *shared_types.OHLCVals) {
	ohlcMap.RLock()
	if list, found := ohlcMap.GetData()[resp.GetMetaData()]; found { // if channelID already exists in the map, then...

		var tmpStartTime time.Time = resp.OHLCArray.EndTime.Add(-time.Minute * time.Duration(sm.GetInterval(resp.GetMetaData())))
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
			list.GetList().AddCandle(&tmp, &types.WSCandles{}, sm.GetInterval(resp.GetMetaData()))
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
}
