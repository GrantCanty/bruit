package settings

type BruitSettings interface {
	Init()
	Wait()
	Add(i int)
	Done()
	CtxDone() <-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
	Load()
	GetBaseCurrency() string
}
