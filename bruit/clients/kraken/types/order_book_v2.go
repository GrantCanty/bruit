package types

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
)

type BaseRespV2WS struct {
	Channel string `json:"channel"`
	Type    string `json:"type"`
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
	Symbol   string       `json:"symbol"`
	Bids     []LevelsV2WS `json:"bids"`
	Asks     []LevelsV2WS `json:"asks"`
	Checksum uint32       `json:"checksum"`
}

type BookRespV2SnapshotJSON struct {
	Symbol   string       `json:"symbol"`
	Bids     *treemap.Map `json:"bids"`
	Asks     *treemap.Map `json:"asks"`
	Checksum uint32       `json:"checksum"`
}

type BookRespV2Update struct {
	BookRespV2Snapshot
	Timestamp time.Time `json:"timestamp"`
}

type BookRespV2UpdateJSON struct {
	BookRespV2SnapshotJSON
	Timestamp time.Time `json:"timestamp"`
}

type LevelsV2WS struct {
	Price    NumericString `json:"price"`
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

type OrderBookWithMutex struct {
	Book  *BookRespV2Update
	Mutex sync.RWMutex
}

type OrderBookWithMutexTree struct {
	Book  *BookRespV2UpdateJSON
	Mutex sync.RWMutex
}

type BookRespV2Success struct {
	Version      string `json:"version"`
	System       string `json:"system"`
	VersionAPI   string `json:"api_version"`
	ConnectionID int    `json:"connection_id"`
}

type SuccessBookResponseV2WS struct {
	BaseRespV2WS
	Data []BookRespV2Success
}

func NumericStringComparator(a, b interface{}) int {
	// Convert strings to float64 for numeric comparison
	var numA float64
	var errA error
	var numB float64
	var errB error

	aNumStr, okA := a.(NumericString)
	bNumStr, okB := b.(NumericString)

	if okA && okB {
		numA, errA = strconv.ParseFloat(string(aNumStr), 64)
		numB, errB = strconv.ParseFloat(string(bNumStr), 64)
	}
	if b, ok := b.(NumericString); ok {
		numB, errB = strconv.ParseFloat(string(b), 64)
	}

	if errA != nil || errB != nil {
		// Fallback to string comparison if parsing fails
		return utils.StringComparator(string(aNumStr), string(bNumStr))
	}

	// Numeric comparison
	switch {
	case numA < numB:
		return -1
	case numA > numB:
		return 1
	default:
		return 0
	}
}

func DeepCopyOrderBook(original BookRespV2UpdateJSON) BookRespV2UpdateJSON {
	// Create a new book
	copy := BookRespV2UpdateJSON{
		// Copy basic fields directly
		Timestamp: original.Timestamp,
		BookRespV2SnapshotJSON: BookRespV2SnapshotJSON{
			Symbol:   original.Symbol,
			Bids:     nil,
			Asks:     nil,
			Checksum: original.Checksum,
		},

		// Add any other simple fields from your struct
	}

	// Create new treemaps for bids and asks
	copy.Bids = treemap.NewWith(NumericStringComparator) // Assuming you're using a custom comparator
	copy.Bids = treemap.NewWith(func(key interface{}, value interface{}) int {
		return -NumericStringComparator(key, value) // Descending order
	})

	copy.Asks = treemap.NewWith(NumericStringComparator)

	// Copy all bid entries
	for _, k := range original.Bids.Keys() {
		v, _ := original.Bids.Get(k)
		copy.Bids.Put(k, v)
	}

	// Copy all ask entries
	for _, k := range original.Asks.Keys() {
		v, _ := original.Asks.Get(k)
		copy.Asks.Put(k, v)
	}

	// Copy any other maps or complex structures in the order book

	return copy
}
