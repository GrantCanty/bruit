# Bruit
A Go client to interact with the Kraken Crypto Exchange WebSocket & Rest APIs

Bruit aims to create a fully automated & algorithmic trading system that's lightweight and powerful. As this is still in version 0, releases are NOT stable. Caution will be taken to ensure that future releases don't have major/breaking changes

## Rough Roadmap (subject to change)
* Finish functions for Kraken order book
* Implementation of the trading ending and strategies
* Implementing the timeseries database (InfluxDB)
* Implementing new modes & main system commands
* Implementing performance tracking
* Connecting to other exchanges

## Usage
To start:
1. Create Kraken account and get API Keys
2. Create a .env file with fields ```API_KEY``` & ```PRIVATE_KEY```. Insert API keys here

#### Subscribe to OHLC
```
package main

import (
	"bruit/bruit"
	"bruit/bruit/clients/kraken"
	"bruit/bruit/influx"
	"bruit/bruit/shared_types"
)

func main() {
	settings := bruit.Settings{}
	settings.Init()

	db := influx.DB{}
	db.Init()

	ohlcMap := shared_types.OHLCVals{}

	k := &kraken.KrakenClient{}
	k.InitClient(&settings)

	go k.PubDecoder(&settings)
	go k.PubListen(&settings, &ohlcMap, db.GetTradeWriter())

	k.SubscribeToOHLC(&settings, []string{"BTC/USD"}, 5)

	go k.DeferChanClose(&settings)
	settings.Wait()
}

### Conventions
Run decoder and listen functions before subscribing to a stream to minimize the chance that of missing messages 
