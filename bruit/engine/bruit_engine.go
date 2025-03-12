package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

type BruitEngine interface {
	Run(s settings.BruitSettings, c clients.BruitCryptoClient, db *influx.DB)
	Stop()
	Wait(s settings.BruitSettings, c clients.BruitCryptoClient)
}
