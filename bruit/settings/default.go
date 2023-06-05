package settings

type DefaultSettings struct {
	Settings

	ConcurrencySettings ConcurrencySettings
	GlobalSettings      Globals
	env                 map[string]string
}

func NewDefaultSettings(parent Settings) Settings {
	return newDefaults(parent)
}

func newDefaults(parent Settings) Settings {
	return &DefaultSettings{Settings: parent}
}

func (s *DefaultSettings) Init() {
	s.ConcurrencySettings.Init()
	s.GlobalSettings.Init()
}

func (s *DefaultSettings) Wait() {
	<-s.ConcurrencySettings.comms
	s.ConcurrencySettings.cancel()
	s.ConcurrencySettings.Wg.Wait()
}

func (s *DefaultSettings) Add(i int) {
	s.ConcurrencySettings.Wg.Add(i)
}

func (s *DefaultSettings) Done() {
	s.ConcurrencySettings.Wg.Done()
}

func (s *DefaultSettings) CtxDone() <-chan struct{} {
	//var d <-chan struct{}
	//d = make(<-chan struct{})
	//s.ConcurrencySettings.Ctx.Done()
	return s.ConcurrencySettings.Ctx.Done()

	//return d
}

func (s *DefaultSettings) GetLoggingToConsole() bool {
	return s.GetLoggingSettings().isLoggingToConsole
}

func (s *DefaultSettings) GetLoggingSettings() LoggingSettings {
	return s.GlobalSettings.Logging
}
