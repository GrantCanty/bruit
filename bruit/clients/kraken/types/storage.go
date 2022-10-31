package types

import (
	"bruit/bruit/shared_types"
)

type KrakenMetaData struct {
	ChannelID   int
	ChannelName string
	Pair        string
}

func (k KrakenMetaData) GetChannelID() int {
	return k.ChannelID
}

func (k KrakenMetaData) GetChannelName() string {
	return k.ChannelName
}

func (k KrakenMetaData) GetPair() string {
	return k.Pair
}

func (k KrakenMetaData) Found(metaData shared_types.SubscriptionMetaData) bool {
	switch data := metaData.(type) {
	case *KrakenMetaData:
		if data.ChannelID == k.ChannelID && data.ChannelName == k.ChannelName && data.Pair == k.Pair {
			return true
		}
	}
	return false
}

type KrakenOHLCSubscriptionData struct {
	Interval int
	Status   string
}

func (k KrakenOHLCSubscriptionData) GetData() shared_types.SubscriptionData {
	return KrakenOHLCSubscriptionData{k.Interval, k.Status}
}

/*type KrakenOHLCResponseHolder struct {
	ChannelID   int
	Interval    int64
	ChannelName string
	Pair        string
	//MetaData KrakenMetaData
	List shared_types.List
}*/

/*func (o KrakenOHLCResponseHolder) GetChannelID() int {
	return o.ChannelID
}*/

/*func (o *KrakenOHLCResponseHolder) GetList() *shared_types.List {
	return &o.List
}*/

/*func (k KrakenOHLCResponseHolder) GetMetaData() {

}*/

/*func (o KrakenOHLCResponseHolder) GetInterval() int64 {
	return o.Interval
}*/

/*type OHLCVals struct {
	Vals  map[int]shared_types.OhlcResponseHolder
	Mutex sync.RWMutex
}

func (ohlcVals *OHLCVals) Set(key int, data shared_types.OhlcResponseHolder) {
	log.Println("-- Setting up map --")
	if ohlcVals.Vals == nil {
		ohlcVals.Vals = make(map[int]shared_types.OhlcResponseHolder)
	}
	ohlcVals.Vals[key] = data
}

func (ohlcVals *OHLCVals) RLock() {
	ohlcVals.Mutex.RLock()
}

func (ohlcVals *OHLCVals) RUnlock() {
	ohlcVals.Mutex.RUnlock()
}

func (ohlcVals *OHLCVals) Lock() {
	ohlcVals.Mutex.Lock()
}

func (ohlcVals *OHLCVals) Unlock() {
	ohlcVals.Mutex.Unlock()
}

func (ohlcVals *OHLCVals) GetVals() map[int]shared_types.OhlcResponseHolder {
	return ohlcVals.Vals
}

func (ohlcVals *OHLCVals) GetMutex() *sync.RWMutex {
	return &ohlcVals.Mutex
}
*/
