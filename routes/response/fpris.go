// SPDX-License-Identifier: EUPL-1.2

//nolint:tagliatelle
package response

import "git.zzdats.lv/edim/api-wallet/routes/object"

type FprisResponse struct {
	object.FprisData
	IssuanceDate string `json:"issuance_date"`
	ExpiryDate   string `json:"expiry_date"`
}
