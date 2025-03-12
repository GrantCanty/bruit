package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

type emptyEngine int

func (e emptyEngine) Run(s settings.BruitSettings, c clients.BruitCryptoClient, db influx.DB) {
	return
}

func (e emptyEngine) Stop() {
	return
}

func (e emptyEngine) Wait(s settings.BruitSettings, c clients.BruitCryptoClient) {
	return
}

/*func Engine() BruitEngine {
	return new(emptyEngine)
}*/
