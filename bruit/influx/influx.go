package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
)

type DB struct {
	Client influxdb2.Client
	Writers
}

type Writers struct {
	tradeWriter api.WriteAPI
}

func (db *DB) Init() {
	db.initClient()
	db.initWriters()
}

func (db *DB) initClient() {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	if key, found := env["INFLUX_KEY"]; found {
		db.Client = influxdb2.NewClient("http://localhost:8086", key)
	} else {
		panic("INFLUX_KEY not found")
	}
}

func (db *DB) initWriters() {
	env, err := godotenv.Read()
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
