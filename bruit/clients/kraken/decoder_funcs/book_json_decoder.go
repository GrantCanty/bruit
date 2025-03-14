package decoders

import (
	"bruit/bruit/clients/kraken/types"
	"bytes"
	"encoding/json"
	"hash/crc32"
	"log"
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
)

// make this run in parallel
func verifyLevelTree(resp treemap.Map, strBuilder *strings.Builder) {
	var priceNum types.NumericString
	var qtyNum types.NumericString
	var ok bool

	var price string
	var qty string

	respItt := resp.Iterator()
	for respItt.Begin(); respItt.Next(); {
		priceNum, ok = respItt.Key().(types.NumericString)
		if !ok {
			log.Printf("Expected float64 key, got %T\n %d", respItt.Key(), respItt.Key())
			continue
		}

		qtyNum, ok = respItt.Value().(types.NumericString)
		if !ok {
			log.Printf("Expected float64 key, got %T\n %d", respItt.Value(), respItt.Value())
			continue
		}

		price = strings.Replace(string(priceNum), ".", "", -1)
		qty = strings.Replace(string(qtyNum), ".", "", -1)

		price = strings.TrimLeft(price, "0")
		qty = strings.TrimLeft(qty, "0")

		strBuilder.WriteString(price)
		strBuilder.WriteString(qty)
	}
}

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

		price = strings.Replace(string(priceNum), ".", "", -1)
		qty = strings.Replace(string(qtyNum), ".", "", -1)

		price = strings.TrimLeft(price, "0")
		qty = strings.TrimLeft(qty, "0")

		strBuilder.WriteString(price)
		strBuilder.WriteString(qty)
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

func VerifyChecksumUpdate(book types.BookRespV2UpdateJSON, resp types.UpdateBookRespV2WS) bool {
	crc32q := crc32.MakeTable(crc32.IEEE)

	var priceAsks strings.Builder
	var priceBids strings.Builder

	verifyLevelTree(*book.Asks, &priceAsks)
	verifyLevelTree(*book.Bids, &priceBids)

	priceAsks.WriteString(priceBids.String())
	cs := crc32.Checksum([]byte(priceAsks.String()), crc32q)
	return cs == resp.Data[0].Checksum
}

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
