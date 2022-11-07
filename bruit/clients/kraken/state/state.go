package state

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
)

/**
 * should have a client manager with managers for subscriptions(channelIDs, intervals for OHLC, status) state(balances, ),
**/

type StateManager struct {
	Account AccountManager
	Client  ClientManager
}

type AccountManager struct {
	balancesWithStaking    types.AccountBalanceResp
	balancesWithoutStaking types.AccountBalanceResp
}

type ClientManager struct {
	subscriptions map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData
}
