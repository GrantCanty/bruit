package rest

import (
	kraken_data "bruit_new/bruit/clients/kraken/client_data"
	"bruit_new/bruit/clients/kraken/types"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type RestClient struct {
	client http.Client
}

func (c *RestClient) PublicRequest(url_path string, values url.Values, returnType interface{}) (interface{}, error) {
	resp, err := c.doRequest(url_path, values, nil, returnType)

	return resp, err
}

func (c *RestClient) PrivateRequest(url_path string, values url.Values, key string, secretKey string, returnType interface{}) (interface{}, error) {
	ssecret, _ := base64.StdEncoding.DecodeString(secretKey)

	// Create signature
	signature := createSignature(url_path, values, ssecret)
	fullPath := strings.Join([]string{kraken_data.RestUrl, url_path}, "")

	// Add Key and signature to request headers
	headers := map[string]string{
		"API-Key":  key,
		"API-Sign": signature,
	}

	resp, err := c.doRequest(fullPath, values, headers, returnType)

	return resp, err
}

func (c *RestClient) doRequest(url_path string, values url.Values, headers map[string]string, returnType interface{}) (interface{}, error) {
	// Create request
	req, err := http.NewRequest("POST", url_path, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #1 (%s)", err.Error())
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! #2 (%s)", err.Error())
	}
	defer resp.Body.Close()

	// Read request
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
		return nil, fmt.Errorf("Could not execute request! #6 (%s)", err.Error())
	}

	// Check for Kraken API error
	if len(jsonData.Error) > 0 {
		return nil, fmt.Errorf("Could not execute request! #7 (%s)", jsonData.Error)
	}

	return jsonData.Result, nil
}
