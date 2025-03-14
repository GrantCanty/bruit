package state

import (
	"bruit/bruit/clients/kraken/types"
	"strings"

	"github.com/shopspring/decimal"
)

func (am *AccountManager) initAccount(bals types.AccountBalanceResp) {
	DeleteZeroBalances(bals)
	am.balancesWithStaking = bals

	newBals := CopyBals(bals)
	ProcessStakingBalances(newBals)
	am.balancesWithoutStaking = newBals
}

func (am *AccountManager) GetBalancesWithStaking() types.AccountBalanceResp {
	am.mutex.Lock()
	tmp := am.balancesWithStaking
	am.mutex.Unlock()

	return tmp
}

func (am *AccountManager) GetBalancesWithoutStaking() types.AccountBalanceResp {
	am.mutex.Lock()
	tmp := am.balancesWithoutStaking
	am.mutex.Unlock()

	return tmp
}

func CopyMap(bals types.AccountBalanceResp) types.AccountBalanceResp {
	//var newMap types.AccountBalanceResp
	newMap := make(types.AccountBalanceResp)
	for pair, amount := range bals {
		newMap[pair] = amount
	}
	return newMap
}

func DeleteZeroBalances(bals types.AccountBalanceResp) {
	zero := decimal.New(0, 0)
	for pair, amount := range bals {
		if !amount.GreaterThan(zero) {
			delete(bals, pair)
		}
	}
}

func ProcessStakingBalances(bals types.AccountBalanceResp) {
	for pair, amount := range bals {
		if strings.Contains(pair, ".S") {
			unStakedPair := strings.Split(pair, ".S")[0]
			bals[unStakedPair] = bals[pair].Add(amount)
			delete(bals, pair)
		}
	}
}

func CopyBals(bals types.AccountBalanceResp) types.AccountBalanceResp {
	//var balsCopy types.AccountBalanceResp
	balsCopy := make(types.AccountBalanceResp)
	for pair, amount := range bals {
		balsCopy[pair] = amount
	}
	return balsCopy
}
