package rest

import (
	"bruit/bruit/clients/kraken/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
)

func decode(resp *http.Response, returnType interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #3 (%s)", err.Error())
	}

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("Could not execute request #4! (%s)", err.Error())
	}
	if mimeType != "application/json" {
		return nil, fmt.Errorf("Could not execute request #5! (%s)", fmt.Sprintf("Response Content-Type is '%s', but should be 'application/json'.", mimeType))
	}

	// Parse request
	var jsonData types.RestResp

	// Set the KrakenResponse.Result to typ so `json.Unmarshal` will
	// unmarshal it into given type, instead of `interface{}`.
	if returnType != nil {
		jsonData.Result = returnType
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		log.Println("resp from ws: ", string(body))
		log.Printf("%sv", returnType)
		return nil, fmt.Errorf("Could not execute request! #6 (%s)", err.Error())
	}

	// Check for Kraken API error
	if len(jsonData.Error) > 0 {
		return nil, fmt.Errorf("Could not execute request! #7 (%s)", jsonData.Error)
	}

	return jsonData.Result, nil
}
