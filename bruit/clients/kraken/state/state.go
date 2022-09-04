package state

import "github.com/shopspring/decimal"

/**
 * should have a client manager with managers for subscriptions(channelIDs, intervals for OHLC, status) state(balances, ),
**/

type Manager struct {
	State  StateManager
	Client ClientManager
}

type StateManager struct {
	Balances map[string]decimal.Decimal
}

type ClientManager struct {
}

func (m *Manager) InitManager() {
	m.State.InitState()
	m.Client.InitClient()
}

func (s *StateManager) InitState() {
	s.Balances = make(map[string]decimal.Decimal)
}

func (c *ClientManager) InitClient() {

}
