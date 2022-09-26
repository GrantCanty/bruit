package types

import (
	"bruit_new/bruit/shared_types"
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

type OHLCResponse struct {
	ChannelID   int
	OHLCArray   WSCandles
	ChannelName string
	Pair        string
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
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
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
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
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
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
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
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
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
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (c *WSCandles) NewStartTime(newTime time.Time) {
	c.StartTime.Time = newTime
}

func (ws WSCandles) GetCandle() shared_types.Candle {
	return &ws
}

func (ws WSCandles) GetStartTime() shared_types.UnixTime {
	return ws.StartTime
}

func (ws WSCandles) GetEndTime() shared_types.UnixTime {
	return ws.EndTime
}

func (ws WSCandles) GetHigh() decimal.Decimal {
	return ws.High
}

func (ws *WSCandles) SetHigh(num decimal.Decimal) {
	ws.High = num
}

func (ws WSCandles) GetLow() decimal.Decimal {
	return ws.Low
}

func (ws *WSCandles) SetLow(num decimal.Decimal) {
	ws.Low = num
}

func (ws WSCandles) GetClose() decimal.Decimal {
	return ws.Close
}

func (ws *WSCandles) SetClose(num decimal.Decimal) {
	ws.Close = num
}

func (ws WSCandles) GetVWAP() decimal.Decimal {
	return ws.VWAP
}

func (ws *WSCandles) SetVWAP(num decimal.Decimal) {
	ws.VWAP = num
}
func (ws WSCandles) GetVolume() decimal.Decimal {
	return ws.Volume
}

func (ws *WSCandles) SetVolume(num decimal.Decimal) {
	ws.Volume = num
}

func (ws WSCandles) GetCount() int {
	return ws.Count
}

func (ws *WSCandles) SetCount(num int) {
	ws.Count = num
}
