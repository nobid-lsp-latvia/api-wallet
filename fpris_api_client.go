// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"fmt"
	"strings"

	"git.zzdats.lv/edim/api-wallet/routes/response"

	"azugo.io/azugo"
)

type FprisAPIClient struct {
	url string
}

func NewFprisAPIClient(url string) *FprisAPIClient {
	return &FprisAPIClient{
		url: strings.TrimSuffix(url, "/"),
	}
}

func (c *FprisAPIClient) GetFprisData(ctx *azugo.Context) (*response.FprisResponse, error) {
	client := ctx.HTTPClient().WithBaseURL(c.url)

	fprisRes := &response.FprisResponse{}

	err := client.GetJSON("/1.0/pid", fprisRes, ctx.Header.InheritAuthorization())
	if err != nil {
		return nil, fmt.Errorf("failed to call FPRIS: %w", err)
	}

	return fprisRes, nil
}
