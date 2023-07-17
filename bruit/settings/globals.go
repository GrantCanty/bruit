package settings

type Globals struct {
	RunTime      RunTimeSettings
	Logging      LoggingSettings
	BaseCurrency string
	RiskLevel    RiskSettings
	BookDepth    int
}

func (g *Globals) Init(runTimes map[string]bool, logging map[string]bool) {
	g.RunTime.init(runTimes)
	g.Logging.init(logging)
	g.BaseCurrency = "USD"
	g.BookDepth = 100
	g.RiskLevel.init()
}

func (g Globals) GetBaseCurrency() string {
	return g.BaseCurrency
}

type RunTimeSettings struct {
	IsSystemsTesting    bool
	IsBackTesting       bool
	IsPaperTrading      bool
	IsProductionTrading bool
}

func (r *RunTimeSettings) init(m map[string]bool) {
	r.IsSystemsTesting = m["ISSYSTEMSTESTING"]
	r.IsBackTesting = m["ISBACKTESTING"]
	r.IsPaperTrading = m["ISPAPERTRADING"]
	r.IsProductionTrading = m["ISPRODUCTION"]
}

type LoggingSettings struct {
	isLoggingToDB      bool
	isLoggingToConsole bool
}

type RiskSettings struct {
	AssetsTraded TradableAssetSettings
	Risk         RiskLevels
}

type TradableAssetSettings struct {
	TradesCrypto bool
	TradesStocks bool
	TradesForex  bool
}

type RiskLevels struct {
	isLevelOne   bool
	isLevelTwo   bool
	isLevelThree bool
}

func (l *LoggingSettings) init(m map[string]bool) {
	l.isLoggingToDB = m["ISLOGGINGTODB"]
	l.isLoggingToConsole = m["ISLOGGINGTOCONSOLE"]
}

func (l LoggingSettings) GetLoggingDB() bool {
	return l.isLoggingToDB
}

func (l LoggingSettings) GetLoggingConsole() bool {
	return l.isLoggingToConsole
}

func (r *RiskSettings) init() {
	r.AssetsTraded.init()
	r.Risk.init()
}

func (a *TradableAssetSettings) init() {
	a.TradesCrypto = true
	a.TradesForex = false
	a.TradesStocks = false
}

func (l *RiskLevels) init() {
	l.isLevelOne = false
	l.isLevelTwo = false
	l.isLevelThree = true
}
