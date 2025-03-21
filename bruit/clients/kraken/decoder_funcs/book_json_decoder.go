package decoders

import (
	"bruit/bruit/clients/kraken/types"
	"bytes"
	"encoding/json"
	"log"
)

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

	if ok := types.VerifyChecksumSnapshot(book); ok && testing {
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

func StatusBookResponseV2WS(byteResponse []byte, testing bool) (*types.StatusBookResponseV2WS, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing {
		log.Println("in SuccessBookResponseDecoverV2 func")
	}

	// decodes byteResponse
	var conn types.StatusBookResponseV2WS
	err := decoder.Decode(&conn)
	if err != nil {
		if testing {
			log.Println("SuccessBookResponseDecoverV2 error: ", err)
		}
		return nil, err
	}

	return &conn, err
}

func SubscribeResponseV2WS(byteResponse []byte, testing bool) (*types.SubscribeResponseV2WS, error) {
	reader := bytes.NewReader(byteResponse)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if testing {
		log.Println("in SubscribeResponseV2WS func")
	}

	// decodes byteResponse
	var conn types.SubscribeResponseV2WS
	err := decoder.Decode(&conn)
	if err != nil {
		if testing {
			log.Println("SubscribeResponseV2WS error: ", err)
		}
		return nil, err
	}

	return &conn, err
}
