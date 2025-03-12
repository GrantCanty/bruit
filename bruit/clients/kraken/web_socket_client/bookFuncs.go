package web_socket

/*import (
	"bruit/bruit/clients/kraken/types"
	"log"
	"sort"

	"github.com/shopspring/decimal"
)

func remove(slice []types.Level, s int) []types.Level {
	return append(slice[:s], slice[s+1:]...)
}

func RemovePriceFromBids(bids []types.Level, price decimal.Decimal) []types.Level {
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.LessThanOrEqual(price) })
	if i < len(bids) && bids[i].Price.Equals(price) {
		return remove(bids, i)
	} else {
		return bids
	}
}

func InsertPriceInBids(bids []types.Level, price decimal.Decimal, volume decimal.Decimal) []types.Level {
	bids = RemovePriceFromBids(bids, price)
	level := types.Level{Price: price, Volume: volume}
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.GreaterThanOrEqual(price) })
	bids = append(bids, types.Level{})
	copy(bids[i+1:], bids[i:])
	bids[i] = level
	return bids
}

func RemovePriceFromAsks(asks []types.Level, price decimal.Decimal) []types.Level {
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThanOrEqual(price) })
	if i < len(asks) && asks[i].Price.Equals(price) {
		return remove(asks, i)
	} else {
		return asks
	}
}

func InsertPriceInAsks(asks []types.Level, price decimal.Decimal, volume decimal.Decimal) []types.Level {
	asks = RemovePriceFromAsks(asks, price)
	level := types.Level{Price: price, Volume: volume}
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThan(price) })
	asks = append(asks, types.Level{})
	copy(asks[i+1:], asks[i:])
	asks[i] = level
	return asks

}

func CreateInitial(book map[string]interface{}, key string) []types.Level {
	var list []types.Level = make([]types.Level, 0)
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
		list = append(list, types.Level{Price: price, Volume: vol})
	}
	return list
}
*/