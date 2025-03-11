package kraken_data

import "strings"

var (
	Priv_ws_token string
)

const (
	PubWSUrl  = "wss://ws.kraken.com/"
	PrivWSUrl = "wss://ws-auth.kraken.com/"
	V2WSUrl = "wss://ws.kraken.com/v2"

	RestUrl = "https://api.kraken.com"

	PubRestUrl    = "/0/public"
	OHLCUrl       = "/OHLC"
	AssetPairsUrl = "/AssetPairs"
	AssetsUrl     = "/Assets"

	PrivRestUrl      = "/0/private"
	AccountBalancUrl = "/Balance"
	WSTokenUrl       = "/GetWebSocketsToken"
)

func GetOHLCIntervals() []int {
	return []int{1, 5, 15, 30, 60, 240, 1440, 10080, 21600}
}

func GetBookDepthLevels() []int {
	return []int{10,25,100,500,1000}
}

func GetPubWSUrl() string {
	return PubWSUrl
}

func GetPrivWSUrl() string {
	return PrivWSUrl
}

func GetV2WsURL() string {
	return V2WSUrl
}

func GetOHLCUrl() string {
	return strings.Join([]string{RestUrl, PrivRestUrl, WSTokenUrl}, "")
}
