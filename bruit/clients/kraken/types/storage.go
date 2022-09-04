package types

import (
	"bruit/bruit/shared_types"
	"log"
	"sync"
)

type KrakenOHLCResponseHolder struct {
	ChannelID   int
	Interval    int64
	ChannelName string
	Pair        string
	List        shared_types.List
}

func (o KrakenOHLCResponseHolder) GetChannelID() int {
	return o.ChannelID
}

func (o *KrakenOHLCResponseHolder) GetList() *shared_types.List {
	return &o.List
}

func (o KrakenOHLCResponseHolder) GetInterval() int64 {
	return o.Interval
}

type OHLCVals struct {
	Vals  map[int]*shared_types.OhlcResponseHolder
	Mutex sync.RWMutex
}

func (ohlcVals *OHLCVals) Set(key int, data *shared_types.OhlcResponseHolder) {
	log.Println("-- Setting up map --")
	if ohlcVals.Vals == nil {
		ohlcVals.Vals = make(map[int]*shared_types.OhlcResponseHolder)
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

func (ohlcVals *OHLCVals) GetVals() map[int]*shared_types.OhlcResponseHolder {
	return ohlcVals.Vals
}

func (ohlcVals *OHLCVals) GetMutex() *sync.RWMutex {
	return &ohlcVals.Mutex
}
