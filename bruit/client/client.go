package client

import (
	"bruit/bruit/clients"
	"bruit/bruit/engine"
	"bruit/bruit/influx"
	"bruit/bruit/settings"
	"log"
)

type BruitClient interface {
	Init(cli clients.BruitCryptoClient)
	Run()
	Wait()
	Stop()
}

type Client struct {
	CryptoClient clients.BruitCryptoClient
	Engine       engine.BruitEngine
	Settings     settings.BruitSettings
	DB           influx.DB
}

func (c *Client) Init(cli clients.BruitCryptoClient) {
	c.Settings = settings.NewDefaultSettings(c.Settings)
	c.Settings.InitSettings()

	switch {
	case c.Settings.IsProduction():
		if c.Settings.GetLoggingToConsole() {
			log.Println("building production engine")
		}
		c.Engine = engine.NewProductionEngine(c.Engine)
		break
	case c.Settings.IsBackTesting():
		if c.Settings.GetLoggingToConsole() {
			log.Println("building backtesting engine")
		}
		c.Engine = engine.NewBackTestEngine(c.Engine)
		break
	case c.Settings.IsPaperTrading():
		if c.Settings.GetLoggingToConsole() {
			log.Println("building paper trading engine")
		}
		c.Engine = engine.NewPaperTradingEngine(c.Engine)
		break
	case c.Settings.IsSystemsTesting():
		if c.Settings.GetLoggingToConsole() {
			log.Println("building systems testing engine")
		}
		c.Engine = engine.NewSystemsTestingEngine(c.Engine)
		break
	default:
		log.Println("no runtime found")
	}

	c.CryptoClient = cli
	c.CryptoClient.InitClient(c.Settings)

	c.DB = influx.DB{}

	return
}

func (c *Client) Run() {
	go c.Engine.Run(c.Settings, c.CryptoClient, &c.DB)
	return
}

func (c *Client) Wait() {
	c.Engine.Wait(c.Settings, c.CryptoClient)
	return
}

func (c *Client) Stop() {
	c.Engine.Stop()
	return
}
