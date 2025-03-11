package types

import (
	"time"
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
	Checksum int `json:"checksum"`
}

type BookRespV2Update struct {
	Symbol string `json:"symbol"`
	Bids []LevelsV2WS `json:"bids"`
	Asks []LevelsV2WS `json:"asks"`
	Checksum int `json:"checksum"`
	Timestamp time.Time `json:"timestamp"`
}

type LevelsV2WS struct {
	Price float64 `json:"price"`
	Quantity float64 `json:"qty"` 
}