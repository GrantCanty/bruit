package settings

type Settings interface {
	Init()
	Wait()
	Add(i int)
	Done()
	CtxDone() <-chan struct{}
	GetLoggingToConsole() bool
	GetLoggingSettings() LoggingSettings
}

type emptySettings int

func (e emptySettings) Init()
func (e emptySettings) Wait()
func (e emptySettings) Add(i int)
func (e emptySettings) Done()
func (e emptySettings) CtxDone() <-chan struct{}
func (e emptySettings) GetLoggingToConsole() bool
func (e emptySettings) GetLoggingSettings() LoggingSettings

func NewSettings() Settings {
	return new(emptySettings)
}
