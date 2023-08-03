package settings

type emptySettings struct {
}

func NewEmptySettings() BruitSettings {
	return newEmpty()
}

func newEmpty() BruitSettings {
	return &emptySettings{}
}

func (e emptySettings) InitSettings() {
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

func (e *emptySettings) GetBaseCurrency() string {
	return ""
}

func NewSettings() BruitSettings {
	return new(emptySettings)
}
