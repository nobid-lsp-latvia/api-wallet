// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"fmt"
	"strings"

	"git.zzdats.lv/edim/api-wallet/routes/response"

	"azugo.io/azugo"
)

type MDLAPIClient struct {
	url string
}

func NewMDLAPIClient(url string) *MDLAPIClient {
	return &MDLAPIClient{
		url: strings.TrimSuffix(url, "/"),
	}
}

func (c *MDLAPIClient) GetMDLData(ctx *azugo.Context) (*response.MDLResponse, error) {
	client := ctx.HTTPClient().WithBaseURL(c.url)

	mdlRes := &response.MDLResponse{}

	err := client.GetJSON("/1.0/mdl", mdlRes, ctx.Header.InheritAuthorization())
	if err != nil {
		return nil, fmt.Errorf("failed to call MDL: %w", err)
	}

	return mdlRes, nil
}
