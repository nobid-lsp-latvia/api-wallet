// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"strings"

	"git.zzdats.lv/edim/api-wallet/routes/response"

	"azugo.io/azugo"
)

type RTUAPIClient struct {
	url string
}

func NewRTUAPIClient(url string) *RTUAPIClient {
	return &RTUAPIClient{
		url: strings.TrimSuffix(url, "/"),
	}
}

func (c *RTUAPIClient) GetRTUData(ctx *azugo.Context) (*response.RTUResponse, error) {
	client := ctx.HTTPClient().WithBaseURL(c.url)

	rtuRes := &response.RTUResponse{}

	err := client.GetJSON("/1.0/diploma", rtuRes, ctx.Header.InheritAuthorization())
	if err != nil {
		return nil, err
	}

	return rtuRes, nil
}
