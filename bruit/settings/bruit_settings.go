package settings

type BruitSettings interface {
	InitSettings()
	Wait()
	Add(i int)
	Done()
	CtxDone() <-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
	Load()
	GetBaseCurrency() string
	IsProduction() bool
	IsBackTesting() bool
	IsPaperTrading() bool
	IsSystemsTesting() bool
}
