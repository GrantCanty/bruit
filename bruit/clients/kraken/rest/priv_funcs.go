package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"net/url"
	"strconv"
	"time"
)

func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

func getHMacSha512(message, secret []byte) []byte {
	//log.Println("hmac512 params: ", message, secret)
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	//log.Println("shaSum: ", shaSum)
	//log.Println("urlPath: ", string(urlPath))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	//log.Println("macSum: ", macSum)
	return base64.StdEncoding.EncodeToString(macSum)
}

func ReturnNonceValues() url.Values {
	nonceTime := strconv.FormatInt((time.Now().UnixMicro() / 1000), 10)
	//nonceTime := strconv.FormatInt((time.Date(2022, time.August, 8, 0, 0, 0, 0, time.UTC).UnixMicro() / 1000), 10)

	//log.Println("noncetime: ", nonceTime)

	params := url.Values{}
	params.Add("nonce", nonceTime)

	return params
}
