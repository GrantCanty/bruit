package influx

import (
	"bruit/bruit/env"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DB struct {
	Client influxdb2.Client
	Writers
}

type Writers struct {
	tradeWriter api.WriteAPI
}

func (db *DB) InitDB() {
	db.initClient()
	db.initWriters()
}

func (db *DB) initClient() {
	env, err := env.Read("DB")
	if err != nil {
		log.Println("err")
		log.Println(env, err)
		panic(err)
	}

	if key, found := env["INFLUX_KEY"]; found {
		db.Client = influxdb2.NewClient("http://localhost:8086", key)
	} else {
		panic("INFLUX_KEY not found")
	}
}

func (db *DB) initWriters() {
	env, err := env.Read("DB")
	if err != nil {
		panic(err)
	}
	org, found1 := env["INFLUX_ORG_NAME"]
	tradesBucket, found2 := env["INFLUX_TRADES_BUCKET_NAME"]
	if found1 && found2 == true {
		db.Writers.tradeWriter = db.Client.WriteAPI(org, tradesBucket)
	} else {
		panic("missing org and or trades bukcet in .env file")
	}

	db.Writers.tradeWriter = db.Client.WriteAPI("Vert", "Trades")
}

func (db *DB) GetTradeWriter() api.WriteAPI {
	return db.Writers.tradeWriter
}
