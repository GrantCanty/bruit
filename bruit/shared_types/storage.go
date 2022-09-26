package shared_types

import "sync"

/*type OHLCResponseHolder struct {
	ChannelID   int
	Interval    int64
	ChannelName string
	Pair        string
	List        List
}*/

type OHLCVals struct {
	Data  map[int]OhlcResponseHolder
	Mutex sync.RWMutex
}

/*type KrakenOHLCResponseHolder struct {
	Data SubscriptionMetaData
	List List
}*/

type SubscriptionMetaData interface {
	GetData() SubscriptionMetaData
	GetChannelID() int
	GetInterval() int64
}

func (ohlcVals *OHLCVals) Set(key int, data OhlcResponseHolder) {
	if ohlcVals.Data == nil {
		ohlcVals.Data = make(map[int]OhlcResponseHolder)
	}
	ohlcVals.Data[key] = data
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

func (ohlcVals *OHLCVals) GetData() map[int]OhlcResponseHolder {
	return ohlcVals.Data
}

func (ohlcVals *OHLCVals) GetMutex() *sync.RWMutex {
	return &ohlcVals.Mutex
}

/*func (o KrakenOHLCResponseHolder) GetChannelID() int {
	return o.Data.GetChannelID()
}

func (o *KrakenOHLCResponseHolder) GetList() *List {
	return &o.List
}

func (o KrakenOHLCResponseHolder) GetInterval() int64 {
	return o.Data.GetInterval()
}*/
