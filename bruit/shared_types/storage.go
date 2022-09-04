package shared_types

import "sync"

type OHLCResponseHolder struct {
	ChannelID   int
	Interval    int64
	ChannelName string
	Pair        string
	List        List
}

type OHLCVals struct {
	Vals  map[int]*OHLCResponseHolder
	Mutex sync.RWMutex
}

func (ohlcVals *OHLCVals) Set(key int, data *OHLCResponseHolder) {
	if ohlcVals.Vals == nil {
		ohlcVals.Vals = make(map[int]*OHLCResponseHolder)
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
