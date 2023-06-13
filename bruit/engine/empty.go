package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

type emptyEngine int

func (e emptyEngine) Init(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Run(s settings.BruitSettings, c clients.BruitClient, db influx.DB) {
	return
}

func (e emptyEngine) Stop() {
	return
}

func (e emptyEngine) Wait(s settings.BruitSettings, c clients.BruitClient) {
	return
}

/*func Engine() BruitEngine {
	return new(emptyEngine)
}*/
