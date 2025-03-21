package types

import (
	"bruit/bruit/shared_types"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ServerConnectionStatusResponse struct {
	ConnectionID uint64 `json:"connectionID"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	Version      string `json:"version"`
}

type HeartBeat struct {
	Event string `json:"event"`
}

type TradeResponse struct {
	ChannelID   int
	TradeArray  []TradeDataResponse
	ChannelName string
	Pair        string
}

type TradeDataResponse struct {
	Price     decimal.Decimal
	Volume    decimal.Decimal
	Time      shared_types.UnixTime
	Side      string
	OrderType string
	Misc      string
}

type OHLCSuccessResponse struct {
	ChannelID    int              `json:"channelID"`
	ChannelName  string           `json:"channelName"`
	Event        string           `json:"event"`
	Pair         string           `json:"pair"`
	Status       string           `json:"status"`
	Subscription OHLCSubscription `json:"subscription"`
}

func (o OHLCSuccessResponse) GetMetaData() shared_types.SubscriptionMetaData {
	return KrakenMetaData{ChannelID: o.ChannelID, ChannelName: o.ChannelName, Pair: o.Pair}
}

type OHLCResponse struct {
	ChannelID   int
	OHLCArray   WSCandles
	ChannelName string
	Pair        string
}

func (o OHLCResponse) GetMetaData() shared_types.SubscriptionMetaData {
	return KrakenMetaData{ChannelID: o.ChannelID, ChannelName: o.ChannelName, Pair: o.Pair}
}

type WSCandles struct {
	StartTime shared_types.UnixTime
	EndTime   shared_types.UnixTime
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	VWAP      decimal.Decimal
	Volume    decimal.Decimal
	Count     int
}

type OpenOrdersResponse struct {
	Orders         []map[string]OpenOrders
	OrdersString   string
	OrdersSequence OpenOrdersSequence
}

type OpenOrders struct {
	AVGPrice        string      `json:"avg_price"`
	Cost            string      `json:"cost"`
	Descr           Description `json:"descr"`
	ExpireTime      *string     `json:"expiretm"`
	Fee             string      `json:"fee"`
	LimitPrice      string      `json:"limitprice"`
	Misc            string      `json:"misc"`
	Oflags          string      `json:"oflags"`
	OpenTime        string      `json:"opentm"`
	Refid           *string     `json:"refid"`
	StartTime       *string     `json:"starttm"`
	Status          string      `json:"status"`
	StopPrice       string      `json:"stopprice"`
	TimeInForce     string      `json:"timeinforce"`
	UserRef         int         `json:"userref"`
	Volume          string      `json:"vol"`
	VolumeExecution string      `json:"vol_exec"`
}

type Description struct {
	Close     *string `json:"close"`
	Leverage  *string `json:"leverage"`
	Order     string  `json:"order"`
	OrderType string  `json:"ordertype"`
	Pair      string  `json:"pair"`
	Price     string  `json:"price"`
	Price2    string  `json:"price2"`
	OType     string  `json:"type"`
}

type OpenOrdersSequence struct {
	Sequence int64 `json:"sequence"`
}

type CancelOrderResponse struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
	Event        string `json:"event"`
	Status       string `json:"status"`
	Count        int    `json:"count,omitempty"`
}

func (t *TradeResponse) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&t.ChannelID, &t.TradeArray, &t.ChannelName, &t.Pair}
	var length int = len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	var g int = len(tmp)
	if g != length {
		return fmt.Errorf("lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (t *TradeDataResponse) UnmarshalJSON(d []byte) error {
	//dd := [][]interface{}{}
	tmp := []interface{}{&t.Price, &t.Volume, &t.Time, &t.Side, &t.OrderType, &t.Misc}
	var length int = len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	var g int = len(tmp)
	if g != length {
		return fmt.Errorf("lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (o *OpenOrdersResponse) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&o.Orders, &o.OrdersString, &o.OrdersSequence}
	var length int = len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	var g int = len(tmp)
	if g != length {
		return fmt.Errorf("lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (ohlc *OHLCResponse) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&ohlc.ChannelID, &ohlc.OHLCArray, &ohlc.ChannelName, &ohlc.Pair}
	length := len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		return fmt.Errorf("lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (c *WSCandles) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&c.StartTime, &c.EndTime, &c.Open, &c.High, &c.Low, &c.Close, &c.VWAP, &c.Volume, &c.Count}
	length := len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		return fmt.Errorf("lengths don't match: %d != %d", g, length)
	}
	return nil
}

// mOST SETTERS NEED A GUARD CLAUSE TO MAKE SURE THAT THE VALUE IS HIGHER THAN THE PREVIOUS RESPONSE (VOL, COUNT, HIGH, LOW(INVERSE), ). THIS HELPS WITH OUT OF ORDER DATA

func (ws WSCandles) GetCandle() shared_types.Candle {
	return &ws
}

func (ws *WSCandles) SetCandle(data ...interface{}) {
	if len(data) != 9 {
		return
	} else {
		ws.StartTime = data[0].(shared_types.UnixTime)
		ws.EndTime = data[1].(shared_types.UnixTime)
		ws.Open = data[2].(decimal.Decimal)
		ws.High = data[3].(decimal.Decimal)
		ws.Low = data[4].(decimal.Decimal)
		ws.Close = data[5].(decimal.Decimal)
		ws.VWAP = data[6].(decimal.Decimal)
		ws.Volume = data[7].(decimal.Decimal)
		ws.Count = data[8].(int)
	}
}

func (ws WSCandles) GetStartTime() shared_types.UnixTime {
	return ws.StartTime
}

func (ws *WSCandles) SetStartTime(newTime time.Time) {
	ws.StartTime.Time = newTime
}

func (ws WSCandles) GetEndTime() shared_types.UnixTime {
	return ws.EndTime
}

func (ws *WSCandles) SetEndTime(newTime time.Time) {
	ws.EndTime.Time = newTime
}

func (ws WSCandles) GetHigh() decimal.Decimal {
	return ws.High
}

func (ws *WSCandles) SetHigh(num decimal.Decimal) {
	if num.GreaterThan(ws.High) {
		ws.High = num
	}
}

func (ws WSCandles) GetLow() decimal.Decimal {
	return ws.Low
}

func (ws *WSCandles) SetLow(num decimal.Decimal) {
	if num.LessThan(ws.Low) {
		ws.Low = num
	}
}

func (ws WSCandles) GetClose() decimal.Decimal {
	return ws.Close
}

func (ws *WSCandles) SetClose(num decimal.Decimal, vol decimal.Decimal) {
	if vol.GreaterThan(ws.Volume) {
		ws.Close = num
	}
}

func (ws WSCandles) GetVWAP() decimal.Decimal {
	return ws.VWAP
}

func (ws *WSCandles) SetVWAP(num decimal.Decimal, vol decimal.Decimal) {
	if vol.GreaterThan(ws.Volume) {
		ws.VWAP = num
	}
}

func (ws WSCandles) GetVolume() decimal.Decimal {
	return ws.Volume
}

func (ws *WSCandles) SetVolume(num decimal.Decimal) {
	if num.GreaterThan(ws.Volume) {
		ws.Volume = num
	}
}

func (ws WSCandles) GetCount() int {
	return ws.Count
}

func (ws *WSCandles) SetCount(num int, vol decimal.Decimal) {
	if vol.GreaterThan(ws.Volume) {
		ws.Count = num
	}
}
