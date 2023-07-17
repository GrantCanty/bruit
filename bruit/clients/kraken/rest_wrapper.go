package kraken

import (
	kraken_data "bruit/bruit/clients/kraken/client_data"
	rest "bruit/bruit/clients/kraken/rest_client"
	"bruit/bruit/clients/kraken/types"
	"strconv"
	"strings"
)

func (c *KrakenClient) GetAssets() (*types.AssetInfoResp, error) {
	url := strings.Join([]string{kraken_data.RestUrl, kraken_data.PubRestUrl, kraken_data.AssetsUrl}, "")
	resp, err := c.Rest.PublicRequest(url, nil, &types.AssetInfoResp{})

	if err != nil {
		return nil, err
	}
	//data, found := (*resp.(*types.AssetResp))["ASD"]
	return resp.(*types.AssetInfoResp), err
}

func (c *KrakenClient) GetAssetPairs() (*types.AssetPairsResp, error) {
	url := strings.Join([]string{kraken_data.RestUrl, kraken_data.PubRestUrl, kraken_data.AssetPairsUrl}, "")
	resp, err := c.Rest.PublicRequest(url, nil, &types.AssetPairsResp{})

	if err != nil {
		return nil, err
	}

	return resp.(*types.AssetPairsResp), err
}

func (c *KrakenClient) GetOHLC(pair string, interval int) (*types.OHLCResp, error) {
	url := strings.Join([]string{kraken_data.RestUrl, kraken_data.PubRestUrl, kraken_data.OHLCUrl, "?pair=", pair, "&interval=", strconv.Itoa(interval)}, "")
	resp, err := c.Rest.PublicRequest(url, nil, &types.OHLCResp{})

	if err != nil {
		return nil, err
	}

	return resp.(*types.OHLCResp), err
}

// PRIV FUNCS

func (c *KrakenClient) GetAccountBalances() (*types.AccountBalanceResp, error) {
	nonceParams := rest.ReturnNonceValues()
	url := strings.Join([]string{kraken_data.PrivRestUrl, kraken_data.AccountBalancUrl}, "")
	resp, err := c.Rest.PrivateRequest(url, nonceParams, kraken_data.ApiKey, kraken_data.PrivateKey, &types.AccountBalanceResp{})

	if err != nil {
		return nil, err
	}

	return resp.(*types.AccountBalanceResp), err
}

func (c *KrakenClient) GetPrivateWebSokcetKey() (*types.PrivWSKeyResp, error) {
	nonceParams := rest.ReturnNonceValues()
	url := strings.Join([]string{kraken_data.PrivRestUrl, kraken_data.WSTokenUrl}, "")
	resp, err := c.Rest.PrivateRequest(url, nonceParams, kraken_data.ApiKey, kraken_data.PrivateKey, &types.PrivWSKeyResp{})

	if err != nil {
		return nil, err
	}

	return resp.(*types.PrivWSKeyResp), err
}
