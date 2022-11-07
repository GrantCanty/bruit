package state

import "bruit/bruit/shared_types"

func (cm *ClientManager) initClient() {
	cm.subscriptions = make(map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData)
}

func (cm ClientManager) GetSubscriptions() map[shared_types.SubscriptionMetaData]shared_types.SubscriptionData {
	return cm.subscriptions
}

func (cm *ClientManager) AddSubscription(metaData shared_types.SubscriptionMetaData, subData shared_types.SubscriptionData) {
	cm.subscriptions[metaData] = subData
}
