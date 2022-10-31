package shared_types

import "sync"

type OHLCVals struct {
	Data  map[SubscriptionMetaData]*List
	Mutex sync.RWMutex
}

func (ohlcVals *OHLCVals) Set(metaData SubscriptionMetaData, data *List) {
	if ohlcVals.Data == nil {
		ohlcVals.Data = make(map[SubscriptionMetaData]*List)
	}
	ohlcVals.Data[metaData] = data
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

func (ohlcVals *OHLCVals) GetData() map[SubscriptionMetaData]*List {
	return ohlcVals.Data
}

func (ohlcVals *OHLCVals) GetMutex() *sync.RWMutex {
	return &ohlcVals.Mutex
}
