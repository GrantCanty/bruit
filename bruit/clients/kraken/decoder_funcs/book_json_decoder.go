package decoders

import (
	"bruit/bruit/clients/kraken/types"
	"bytes"
	"encoding/json"
	"log"

	//"fmt"
	"hash/crc32"
	"strings"
	//"strconv"
	//"reflect"
)

// make this run in parallel
func verifyLevel(resp []types.LevelsV2WS, strBuilder *strings.Builder) {
	var level types.LevelsV2WS
	var priceNum types.NumericString
	var qtyNum types.NumericString

	var price string
	var qty string

	for ask := range resp {
		level = resp[ask]
		priceNum = level.Price
		qtyNum = level.Quantity
		//log.Println("og: ", priceNum, qtyNum)

		price = strings.Replace(string(priceNum), ".", "", -1)
		qty = strings.Replace(string(qtyNum), ".", "", -1)
		//log.Println("no decimal: ", price, qty)

		price = strings.TrimLeft(price, "0")
		qty = strings.TrimLeft(qty, "0")
		//log.Println("trimmed left: ", price, qty)

		strBuilder.WriteString(price)
		strBuilder.WriteString(qty)
		//log.Println("priceAsks: ", strBuilder.String())
	}
}

// make verifyLevel run in parallel
func verifyChecksumSnapshot(resp types.SnapshotBookRespV2WS) bool {
	crc32q := crc32.MakeTable(crc32.IEEE)

	var priceAsks strings.Builder
	var priceBids strings.Builder

	verifyLevel(resp.Data[0].Asks, &priceAsks)
	verifyLevel(resp.Data[0].Bids, &priceBids)
	priceAsks.WriteString(priceBids.String())

	cs := crc32.Checksum([]byte(priceAsks.String()), crc32q)

	return cs == resp.Data[0].Checksum
}

/*func InitialBookResponseDecoder(byteResponse []byte, now time.Time, testing bool) (*types.BookDecodedResp, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in initial book response decoder func")
	}

	// decodes byteResponse
	var book types.InitialBookResp
	err := decoder.Decode(&book)
	if err != nil {
		if testing == true {
			log.Println("initialBookResponseDecoder error: ", err)
		}
		return nil, err
	}

	// transfer data from book to ob
	var ob types.BookDecodedResp
	ob.TimeReceived = now
	ob.Asks = book.Levels["as"]
	ob.Bids = book.Levels["bs"]

	return &ob, err
}*/

func SnapshotBookResponseDecoderV2(byteResponse []byte, testing bool) (*types.SnapshotBookRespV2WS, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing {
		log.Println("in SnapshotBookResponseDecoderV2 func")
	}

	var book types.SnapshotBookRespV2WS

	err := decoder.Decode(&book)
	if err != nil {
		if testing {
			log.Println("SnapshotBookResponseDecoderV2 error: ", err)
		}
		return nil, err
	}

	if ok := verifyChecksumSnapshot(book); ok && testing {
		log.Println("checksums match")
	}

	return &book, nil
}

func UpdateBookResponseDecoderV2(byteResponse []byte, testing bool) (*types.UpdateBookRespV2WS, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing {
		log.Println("in UpdateBookResponseDecoderV2 func")
	}

	// decodes byteResponse
	var book types.UpdateBookRespV2WS
	err := decoder.Decode(&book)
	if err != nil {
		if testing {
			log.Println("UpdateBookResponseDecoderV2 error: ", err)
		}
		return nil, err
	}

	//verifyChecksum(book)

	return &book, err
}

func SuccessBookResponseDecoverV2(byteResponse []byte, testing bool) (*types.SuccessBookResponseV2WS, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing {
		log.Println("in SuccessBookResponseDecoverV2 func")
	}

	// decodes byteResponse
	var conn types.SuccessBookResponseV2WS
	err := decoder.Decode(&conn)
	if err != nil {
		if testing {
			log.Println("SuccessBookResponseDecoverV2 error: ", err)
		}
		return nil, err
	}

	return &conn, err
}

/*func IncrementalAskOrBidDecoder(byteResponse []byte, testing bool) (*types.UpdateBookWithAsksOrBidsResp, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in incremental ask or bid decoder func")
	}

	// decodes byteResponse
	var asksOrBids types.UpdateBookWithAsksOrBidsResp

	err := decoder.Decode(&asksOrBids)
	if err != nil {
		if testing == true {
			log.Println("incrementalAskOrBidDecoder error: ", err)
		}
		return nil, err
	}

	return &asksOrBids, nil

}

func IncrementalAskAndBidDecoder(byteResponse []byte, testing bool) (*types.UpdateBookWithAsksAndBidsResp, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in incremental ask and bid decoder func")
	}

	// decodes byteResponse
	var asksAndBids types.UpdateBookWithAsksAndBidsResp
	err := decoder.Decode(&asksAndBids)
	if err != nil {
		if testing == true {
			log.Println("incrementalAskAndBidDecoder erroor: ", err)
		}
		return nil, err
	}

	return &asksAndBids, nil
}*/

/*func BookSubscriptionResponseDecoder(byteResponse []byte, testing bool) (*types.BookSuccessResponse, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing == true {
		log.Println("in ohlc subscription response func")
	}
	var ohlc types.BookSuccessResponse
	err := decoder.Decode(&ohlc)
	if err != nil {
		if testing == true {
			log.Println("ohlc subscription response error", err)
		}
		return nil, err
	}
	return &ohlc, err
}*/

/*var resp []interface{}
var ob types.BookDecodedResp

err := json.Unmarshal(byteResponse, &resp)
if err != nil {
	return nil, err
}

for _, element := range resp[1].(map[string]interface{})["as"].([]interface{}) {
	priceStr := element.([]interface{})[0].(string)
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		return nil, err
	}
	volStr := element.([]interface{})[1].(string)
	vol, err := decimal.NewFromString(volStr)
	if err != nil {
		return nil, err
	}
	ob.Asks = append(ob.Asks, types.Level{Price: price, Volume: vol})
}

for _, element := range resp[1].(map[string]interface{})["bs"].([]interface{}) {
	priceStr := element.([]interface{})[0].(string)
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		return nil, err
	}
	volStr := element.([]interface{})[1].(string)
	vol, err := decimal.NewFromString(volStr)
	if err != nil {
		return nil, err
	}
	ob.Bids = append(ob.Bids, types.Level{Price: price, Volume: vol})
}
return &ob, nil
*/
