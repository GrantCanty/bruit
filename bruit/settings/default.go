package settings

import (
	"bruit/bruit/env"
	"log"
	"strconv"
)

type DefaultSettings struct {
	Settings

	ConcurrencySettings ConcurrencySettings
	GlobalSettings      Globals
	env                 map[string]bool
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
	return s.ConcurrencySettings.Ctx.Done()
}

func (s *DefaultSettings) GetLoggingToConsole() bool {
	return s.GetLoggingSettings().isLoggingToConsole
}

func (s *DefaultSettings) GetLoggingSettings() LoggingSettings {
	return s.GlobalSettings.Logging
}

func (s *DefaultSettings) Load() {
	configs, err := env.Read("CONFIG")
	if err != nil {
		panic(err)
	}

	s.initEnv(configs)

	if s.GetLoggingToConsole() {
		switch true {
		case s.env["ISPRODUCTION"]:
			log.Println("prod")
			break
		case s.env["ISBACKTESTING"]:
			log.Println("back")
			break
		case s.env["ISPAPERTRADING"]:
			log.Println("paper")
			break
		case s.env["ISSYSTEMSTESTING"]:
			log.Println("system")
			break
		}
	}
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}
func (s *DefaultSettings) makeMap() {
	s.env = make(map[string]bool)

	s.env["ISPRODUCTION"] = false
	s.env["ISBACKTESTING"] = false
	s.env["ISPAPERTRADING"] = false
	s.env["ISSYSTEMSTESTING"] = false
}

func (s *DefaultSettings) initEnv(configs map[string]string) {
	s.makeMap()
	keys := getKeys(s.env)
	var err interface{}

	var trueCount, falseCount int = 0, 0
	for i := 0; i < len(s.env); i++ {
		s.env[keys[i]], err = strconv.ParseBool(configs[keys[i]])
		if err != nil {
			panic(err)
		}
		if s.env[keys[i]] == true {
			trueCount++
		} else {
			falseCount++
		}

	}
	if trueCount != 1 || falseCount != 3 {
		panic("Incorrect count of true and false runtimes")
	}
}
