package settings

import (
	"bruit/bruit/env"
	"errors"
	"log"
	"strconv"
)

var ErrBadRunTimeSettings error = errors.New("bad runtime settings - incorrect count of true/false")

type settings struct {
	BruitSettings

	ConcurrencySettings ConcurrencySettings
	GlobalSettings      Globals
	runTimes            map[string]bool
	logging             map[string]bool
}

func NewDefaultSettings(parent BruitSettings) BruitSettings {
	return newDefaults(parent)
}

func newDefaults(parent BruitSettings) BruitSettings {
	return &settings{BruitSettings: parent}
}

func (s *settings) InitSettings() error {
	if err := s.Load(); err != nil {
		return err
	}
	s.ConcurrencySettings.Init()
	s.GlobalSettings.Init(s.runTimes, s.logging)
	return nil
}

func (s *settings) Wait() {
	<-s.ConcurrencySettings.comms
	s.ConcurrencySettings.cancel()
	s.ConcurrencySettings.Wg.Wait()
}

func (s *settings) Add(i int) {
	s.ConcurrencySettings.Wg.Add(i)
}

func (s *settings) Done() {
	s.ConcurrencySettings.Wg.Done()
}

func (s *settings) CtxDone() <-chan struct{} {
	return s.ConcurrencySettings.Ctx.Done()
}

func (s *settings) GetLoggingToConsole() bool {
	return s.GetLoggingSettings().isLoggingToConsole
}

func (s *settings) GetLoggingSettings() LoggingSettings {
	return s.GlobalSettings.Logging
}

func (s *settings) Load() error {
	configs, err := env.Read("CONFIG")
	if err != nil {
		return err
	}

	s.initRunTimes(configs)
	s.initLogging(configs)

	if s.logging["ISLOGGINGTOCONSOLE"] {
		log.Println(s.runTimes)
		log.Println(s.logging)

		switch true {
		case s.runTimes["ISPRODUCTION"]:
			log.Println("prod")
			fallthrough
		case s.runTimes["ISBACKTESTING"]:
			log.Println("back")
			fallthrough
		case s.runTimes["ISPAPERTRADING"]:
			log.Println("paper")
			fallthrough
		case s.runTimes["ISSYSTEMSTESTING"]:
			log.Println("system")
			fallthrough
		case s.logging["ISLOGGINGTOCONSOLE"]:
			log.Println("console")
			fallthrough
		case s.logging["ISLOGGINGTODB"]:
			log.Println("db")
		}
	}
	return nil
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

func (s *settings) makeRunTimes() {
	s.runTimes = make(map[string]bool)

	s.runTimes["ISPRODUCTION"] = false
	s.runTimes["ISBACKTESTING"] = false
	s.runTimes["ISPAPERTRADING"] = false
	s.runTimes["ISSYSTEMSTESTING"] = false
}

func (s *settings) initRunTimes(configs map[string]string) error {
	s.makeRunTimes()
	keys := getKeys(s.runTimes)
	var err error

	var trueCount, falseCount int = 0, 0
	for i := 0; i < len(s.runTimes); i++ {
		s.runTimes[keys[i]], err = strconv.ParseBool(configs[keys[i]])
		if err != nil {
			return err
		}
		if s.runTimes[keys[i]] == true {
			trueCount++
		} else {
			falseCount++
		}

	}
	if trueCount != 1 || falseCount != 3 {
		return ErrBadRunTimeSettings
	}
	return nil
}

func (s *settings) makeLogging() {
	s.logging = make(map[string]bool)

	s.logging["ISLOGGINGTOCONSOLE"] = false
	s.logging["ISLOGGINGTODB"] = false
}

func (s *settings) initLogging(configs map[string]string) error {
	s.makeLogging()
	keys := getKeys(s.logging)
	var err error

	for i := 0; i < len(s.logging); i++ {
		s.logging[keys[i]], err = strconv.ParseBool(configs[keys[i]])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *settings) GetBaseCurrency() string {
	return s.GlobalSettings.GetBaseCurrency()
}

func (s *settings) IsProduction() bool {
	return s.runTimes["ISPRODUCTION"] && !s.runTimes["ISBACKTESTING"] && !s.runTimes["ISPAPERTRADING"] && !s.runTimes["ISSYSTEMSTESTING"]
}

func (s *settings) IsBackTesting() bool {
	return !s.runTimes["ISPRODUCTION"] && s.runTimes["ISBACKTESTING"] && !s.runTimes["ISPAPERTRADING"] && !s.runTimes["ISSYSTEMSTESTING"]
}

func (s *settings) IsPaperTrading() bool {
	return !s.runTimes["ISPRODUCTION"] && !s.runTimes["ISBACKTESTING"] && s.runTimes["ISPAPERTRADING"] && !s.runTimes["ISSYSTEMSTESTING"]
}

func (s *settings) IsSystemsTesting() bool {
	return !s.runTimes["ISPRODUCTION"] && !s.runTimes["ISBACKTESTING"] && !s.runTimes["ISPAPERTRADING"] && s.runTimes["ISSYSTEMSTESTING"]
}
