package engine

import (
	"bruit/bruit/clients"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
)

type BruitEngine interface {
	Init(s settings.BruitSettings, c clients.BruitClient, db *influx.DB)
	Run( /*s settings.BruitSettings, c clients.BruitClient, db influx.DB*/ )
	Stop()
	Wait( /*s settings.BruitSettings, c clients.BruitClient*/ )
}
