package types

import (
	"time"
	"strings"
	"strconv"
)

type BaseRespV2WS struct {
	Channel string `json:"channel"`
	Type string `json:"type"`

}

type SnapshotBookRespV2WS struct {
	BaseRespV2WS
	Data []BookRespV2Snapshot `json:"data"`
}

type UpdateBookRespV2WS struct {
	BaseRespV2WS
	Data []BookRespV2Update `json:"data"`
}

type BookRespV2Snapshot struct {
	Symbol string `json:"symbol"`
	Bids []LevelsV2WS `json:"bids"`
	Asks []LevelsV2WS `json:"asks"`
	Checksum uint32 `json:"checksum"`
}

type BookRespV2Update struct {
	Symbol string `json:"symbol"`
	Bids []LevelsV2WS `json:"bids"`
	Asks []LevelsV2WS `json:"asks"`
	Checksum uint32 `json:"checksum"`
	Timestamp time.Time `json:"timestamp"`
}

type LevelsV2WS struct {
	Price NumericString `json:"price"`
	Quantity NumericString `json:"qty"` 
}

type NumericString string

func (n *NumericString) UnmarshalJSON(data []byte) error {
    // Strip the quotes and keep as string
    *n = NumericString(strings.Trim(string(data), "\""))
    return nil
}

// Convert to float when needed
func (n NumericString) Float64() (float64, error) {
    return strconv.ParseFloat(string(n), 64)
}