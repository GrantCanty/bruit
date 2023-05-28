package settings

type Globals struct {
	RunTime      RunTimeSettings
	Logging      LoggingSettings
	BaseCurrency string
	RiskLevel    RiskSettings
	BookDepth    int
}

func (g *Globals) Init() {
	g.RunTime.init()
	g.Logging.init()
	g.BaseCurrency = "USD"
	g.BookDepth = 100
	g.RiskLevel.init()
}

type RunTimeSettings struct {
	IsSystemsTesting    bool
	IsBackTesting       bool
	IsPaperTrading      bool
	IsProductionTrading bool
}

func (r *RunTimeSettings) init() {
	r.IsSystemsTesting = true
	r.IsBackTesting = false
	r.IsPaperTrading = false
	r.IsProductionTrading = false
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

/*func (g *ConcurrencySettings) Wait() {
	<-g.comms
	g.cancel()
	g.Wg.Wait()
}*/

/*func (s *Settings) GetEnv() map[string]string {
	return s.env
}*/

func (l *LoggingSettings) init() {
	l.isLoggingToDB = true
	l.isLoggingToConsole = true
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
