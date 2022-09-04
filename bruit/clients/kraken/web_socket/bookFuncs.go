package web_socket

import (
	"log"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

type BookResp struct {
	TimeReceived time.Time
	Bids         []Level
	Asks         []Level
}

type Level struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
}

func remove(slice []Level, s int) []Level {
	return append(slice[:s], slice[s+1:]...)
}

func RemovePriceFromBids(bids []Level, price decimal.Decimal) []Level {
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.LessThanOrEqual(price) })
	if i < len(bids) && bids[i].Price.Equals(price) {
		return remove(bids, i)
	} else {
		return bids
	}
}

func InsertPriceInBids(bids []Level, price decimal.Decimal, volume decimal.Decimal) []Level {
	bids = RemovePriceFromBids(bids, price)
	level := Level{Price: price, Volume: volume}
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.GreaterThanOrEqual(price) })
	bids = append(bids, Level{})
	copy(bids[i+1:], bids[i:])
	bids[i] = level
	return bids
}

func RemovePriceFromAsks(asks []Level, price decimal.Decimal) []Level {
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThanOrEqual(price) })
	if i < len(asks) && asks[i].Price.Equals(price) {
		return remove(asks, i)
	} else {
		return asks
	}
}

func InsertPriceInAsks(asks []Level, price decimal.Decimal, volume decimal.Decimal) []Level {
	asks = RemovePriceFromAsks(asks, price)
	level := Level{Price: price, Volume: volume}
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThan(price) })
	asks = append(asks, Level{})
	copy(asks[i+1:], asks[i:])
	asks[i] = level
	return asks

}

func CreateInitial(book map[string]interface{}, key string) []Level {
	var list []Level = make([]Level, 0)
	for _, element := range book[key].([]interface{}) {
		priceInterface := element.([]interface{})[0]
		priceStr := priceInterface.(string)
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			log.Fatal(err)
		}

		volInterface := element.([]interface{})[1]
		volStr := volInterface.(string)
		vol, err := decimal.NewFromString(volStr)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, Level{Price: price, Volume: vol})
	}
	return list
}
