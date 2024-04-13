package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"bruit/bruit/shared_types"
)

type BruitEngine interface {
	Init(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB, str shared_types.Strategy)
	Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB)
	Stop()
	Wait(s settings.BruitSettings, c clients.BruitCryptoClient)
}
