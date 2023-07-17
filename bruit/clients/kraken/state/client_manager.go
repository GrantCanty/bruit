package state

import (
	"bruit/bruit/clients/kraken/types"
	"bruit/bruit/shared_types"
)

func (cm *ClientManager) initClient() {
	cm.subscriptions = make(map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData)
	cm.assets = make(types.AssetInfoResp)
	cm.assetPairs = make(types.AssetPairsResp)
}

func (cm *ClientManager) initAssets(assets types.AssetInfoResp) {
	cm.assets = assets
}

func (cm *ClientManager) initPairs(pairs types.AssetPairsResp) {
	cm.assetPairs = pairs
}

func (cm ClientManager) GetSubscriptions() map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData {
	return cm.subscriptions
}

func (cm *ClientManager) AddSubscription(metaData shared_types.SubscriptionMetaData, subData shared_types.SubscriptionData) {
	cm.subscriptions[metaData] = subData
}

func (cm ClientManager) GetAssetPairs() types.AssetPairsResp {
	return cm.assetPairs
}
