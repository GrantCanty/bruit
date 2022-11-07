package types

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type RestResp struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

type AssetInfoResp map[string]AssetInfo

type AssetInfo struct {
	Aclass   string `json:"aclass"`
	AltName  string `json:"altname"`
	Decimals int    `json:"decimals"`
	Display  int    `json:"display_decimals"`
}

type AssetPairResp struct {
	AltName        string     `json:"altname"`
	WsName         string     `json:"wsname"`
	AclassBase     string     `json:"aclass_base"`
	Base           string     `json:"base"`
	AclassQuote    string     `json:"aclass_quote"`
	Quote          string     `json:"quote"`
	Lot            string     `json:"lot"`
	PairDecimals   int        `json:"pair_decimals"`
	LotDecimals    int        `json:"lot_decimals"`
	LotMultiplier  int        `json:"lot_multiplier"`
	LeverageBuy    []int      `json:"leverage_buy"`
	LeverageSell   []int      `json:"leverage_sell"`
	Fees           []FeesInfo `json:"fees"`
	MakerFees      []FeesInfo `json:"fees_maker"`
	FeeVolCurrency string     `json:"fee_volume_currency"`
	MarginCall     int        `json:"margin_call"`
	MarginStop     int        `json:"margin_stop"`
	Ordermin       string     `json:"order_min"`
}

type AccountBalanceResp map[string]decimal.Decimal

type FeesInfo struct {
	Vol int
	Pct float64
}

type OHLCResp struct {
	Pair map[string][]RestCandles
	Last float64 `json:"last"`
}

type RestCandles struct {
	Time   float64
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	VWAP   decimal.Decimal
	Volume decimal.Decimal
	Count  int
}

func (f *FeesInfo) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&f.Vol, &f.Pct}
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

func (r *OHLCResp) UnmarshalJSON(d []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(d, &m); err != nil {
		return err
	}
	if last, ok := m["last"]; ok {
		if err := json.Unmarshal(last, &r.Last); err != nil {
			return err
		}
		delete(m, "last")
	}
	r.Pair = make(map[string][]RestCandles, len(m))
	for k, v := range m {
		cc := []RestCandles{}
		if err := json.Unmarshal(v, &cc); err != nil {
			return err
		}
		r.Pair[k] = cc
	}
	return nil
}

func (c *RestCandles) UnmarshalJSON(d []byte) error {
	tmp := []interface{}{&c.Time, &c.Open, &c.High, &c.Low, &c.Close, &c.VWAP, &c.Volume, &c.Count}
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

type PrivWSKeyResp struct {
	Expires int    `json:"expires"`
	Token   string `json:"token"`
}
