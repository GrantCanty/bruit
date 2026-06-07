package settings

type BruitSettings interface {
	InitSettings() error
	Wait()
	Add(i int)
	Done()
	CtxDone() <-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
	Load() error
	GetBaseCurrency() string
	IsProduction() bool
	IsBackTesting() bool
	IsPaperTrading() bool
	IsSystemsTesting() bool
}
