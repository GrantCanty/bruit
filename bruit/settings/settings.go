package settings

type Settings interface {
	Init()
	Wait()
	Add(i int)
	Done()
	CtxDone() <-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
	Load()
}

type emptySettings struct {
}

func NewEmptySettings() Settings {
	return newEmpty()
}

func newEmpty() Settings {
	return &emptySettings{}
}

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

func (e emptySettings) CtxDone() <-chan struct{} {
	return make(<-chan struct{})
}

func (e emptySettings) GetLoggingToConsole() bool {
	return false
}

func (e emptySettings) GetLoggingSettings() LoggingSettings {
	return LoggingSettings{}
}

func (e *emptySettings) Load() {
	return
}

func NewSettings() Settings {
	return new(emptySettings)
}
