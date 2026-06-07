package influx

import (
	"bruit/bruit/env"
	"errors"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var ErrMissingOrgAndTradesBucket = errors.New("Missing org and or trades bucket in .env file")
var ErrMissingInfluxKey = errors.New("Missing Influx key in .env file")

type DB struct {
	Client influxdb2.Client
	Writers
}

type Writers struct {
	tradeWriter api.WriteAPI
}

func (db *DB) InitDB() error {
	if err := db.initClient(); err != nil {
		return err
	}
	if err := db.initWriters(); err != nil {
		return err
	}
	return nil
}

func (db *DB) initClient() error {
	env, err := env.Read("DB")
	if err != nil {
		log.Println("err")
		log.Println(env, err)
		return err
	}

	if key, found := env["INFLUX_KEY"]; found {
		db.Client = influxdb2.NewClient("http://localhost:8086", key)
	} else {
		return ErrMissingInfluxKey
	}
	return nil
}

func (db *DB) initWriters() error {
	env, err := env.Read("DB")
	if err != nil {
		return err
	}
	org, found1 := env["INFLUX_ORG_NAME"]
	tradesBucket, found2 := env["INFLUX_TRADES_BUCKET_NAME"]
	if found1 && found2 == true {
		db.Writers.tradeWriter = db.Client.WriteAPI(org, tradesBucket)
	} else {
		return ErrMissingOrgAndTradesBucket
	}

	db.Writers.tradeWriter = db.Client.WriteAPI("Vert", "Trades")
	return nil
}

func (db *DB) GetTradeWriter() api.WriteAPI {
	return db.Writers.tradeWriter
}
