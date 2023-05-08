package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type InitialBookResp struct {
	ChannelID   int
	Levels      map[string][]Level
	ChannelName string
	Pair        string
}

type UpdateBookWithAsksOrBidsResp struct {
	ChannelID   int
	PriceAndVol map[string]interface{}
	ChannelName string
	Pair        string
}

type UpdateBookWithAsksAndBidsResp struct {
	ChannelID   int
	Asks        map[string]interface{}
	Bids        map[string]interface{}
	ChannelName string
	Pair        string
}

type BookDecodedResp struct {
	TimeReceived time.Time
	Bids         []Level
	Asks         []Level
}

type Level struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
	//TS     shared_types.UnixTime
}

type BookSuccessResponse struct {
	ChannelID    int              `json:"channelID"`
	ChannelName  string           `json:"channelName"`
	Event        string           `json:"event"`
	Pair         string           `json:"pair"`
	Status       string           `json:"status"`
	Subscription BookSubscription `json:"subscription"`
}

func (b *InitialBookResp) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&b.ChannelID, &b.Levels, &b.ChannelName, &b.Pair}
	length := len(tmp)
	if err := json.Unmarshal(d, &tmp); err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (b *UpdateBookWithAsksOrBidsResp) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&b.ChannelID, &b.PriceAndVol, &b.ChannelName, &b.Pair}
	length := len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		fmt.Println(tmp)
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (b *UpdateBookWithAsksAndBidsResp) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&b.ChannelID, &b.Asks, &b.Bids, &b.ChannelName, &b.Pair}
	length := len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		fmt.Println(tmp)
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
	}
	return nil
}

func (l *Level) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&l.Price, &l.Volume}
	length := len(tmp)
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	g := len(tmp)
	if g != length {
		fmt.Println(tmp)
		return fmt.Errorf("Lengths don't match: %d != %d", g, length)
	}
	return nil
}
