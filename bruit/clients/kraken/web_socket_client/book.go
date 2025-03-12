package web_socket

import (
	/*"bruit/bruit/clients/kraken/types"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/shopspring/decimal"*/
)

//orderBook := make(map[string]*types.V2OrderBookMutex)

/*func GetCheckSumInput(bids []types.Level, asks []types.Level) string {
	var str strings.Builder
	for _, level := range asks[:10] {
		price := level.Price.StringFixed(5)
		price = strings.Replace(price, ".", "", 1)
		price = strings.TrimLeft(price, "0")
		str.WriteString(price)

		volume := level.Volume.StringFixed(8)
		volume = strings.Replace(volume, ".", "", 1)
		volume = strings.TrimLeft(volume, "0")
		str.WriteString(price)
	}
	for _, level := range bids[:10] {
		price := level.Price.StringFixed(5)
		price = strings.Replace(price, ".", "", 1)
		price = strings.TrimLeft(price, "0")
		str.WriteString(price)

		volume := level.Volume.StringFixed(8)
		volume = strings.Replace(volume, ".", "", 1)
		volume = strings.TrimLeft(volume, "0")
		str.WriteString(price)
	}
	return str.String()
}

func verifyOrderBookCheckSum(bids []types.Level, asks []types.Level, checkSum string) {
	checSumInput := GetCheckSumInput(bids, asks)
	crc := crc32.ChecksumIEEE([]byte(checSumInput))
	if fmt.Sprint(crc) != checkSum {
		panic(fmt.Sprint("not the same ", crc, " ", checkSum))
	}
}

func getPriceAndVolume(ino interface{}) (decimal.Decimal, decimal.Decimal, int) {
	el := ino.([]interface{})
	price, err := decimal.NewFromString(el[0].(string))
	if err != nil {
		panic(err)
	}
	volume, err := decimal.NewFromString(el[1].(string))
	if err != nil {
		panic(err)
	}
	return price, volume, len(el)
}
*/