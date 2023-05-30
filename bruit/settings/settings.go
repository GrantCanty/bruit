package settings

type Settings interface {
	Init()
	Wait()
	Add(i int)
	Done()
	CtxDone() //<-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
}

type emptySettings int

func (e emptySettings) Init() {
	return
}

func (e emptySettings) Wait() {
	return
}

func (e emptySettings) Add(i int) {
	return
}

func (e emptySettings) Done() {
	return
}

func (e emptySettings) CtxDone() /*<-chan struct{}*/ {
	return //nil
}

func (e emptySettings) GetLoggingToConsole() bool {
	return false
}

func (e emptySettings) GetLoggingSettings() LoggingSettings {
	return LoggingSettings{}
}

func NewSettings() Settings {
	return new(emptySettings)
}
