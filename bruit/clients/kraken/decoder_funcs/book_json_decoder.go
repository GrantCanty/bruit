package decoders

import (
	"bruit/bruit/clients/kraken/types"
	"bytes"
	"encoding/json"
	"log"
	"time"
)

func BookJsonDecoder(response string, testing bool) []interface{} {
	var resp []interface{}
	//byteResponse := []byte(response)
	log.Println(response)

	/*_, err := hbResponseDecoder(byteResponse, testing)
	if err != nil {
		resp, err = rr()
	}*/

	return resp
}

func InitialBookResponseDecoder(byteResponse []byte, now time.Time, testing bool) (*types.BookDecodedResp, error) {
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
}

func IncrementalAskOrBidDecoder(byteResponse []byte, testing bool) (*types.UpdateBookWithAsksOrBidsResp, error) {
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
}

func BookSubscriptionResponseDecoder(byteResponse []byte, testing bool) (*types.BookSuccessResponse, error) {
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
}

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
