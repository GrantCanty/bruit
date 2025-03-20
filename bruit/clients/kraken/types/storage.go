package types

import (
	"bruit/bruit/shared_types"
)

type KrakenMetaData struct {
	ChannelID   int
	ChannelName string
	Pair        string
}

func (k KrakenMetaData) GetChannelID() int {
	return k.ChannelID
}

func (k KrakenMetaData) GetChannelName() string {
	return k.ChannelName
}

func (k KrakenMetaData) GetPair() string {
	return k.Pair
}

func (k KrakenMetaData) Found(metaData shared_types.SubscriptionMetaData) bool {
	switch data := metaData.(type) {
	case *KrakenMetaData:
		if data.ChannelID == k.ChannelID && data.ChannelName == k.ChannelName && data.Pair == k.Pair {
			return true
		}
	}
	return false
}

type KrakenOHLCSubscriptionData struct {
	Interval int
	Status   string
}

func (k KrakenOHLCSubscriptionData) GetData() shared_types.SubscriptionData {
	return KrakenOHLCSubscriptionData{k.Interval, k.Status}
}
