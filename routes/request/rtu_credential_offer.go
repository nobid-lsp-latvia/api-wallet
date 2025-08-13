// SPDX-License-Identifier: EUPL-1.2

package request

import (
	"git.zzdats.lv/edim/api-wallet/routes/object"
	"git.zzdats.lv/edim/api-wallet/routes/response"
)

type RTUCredentialOfferRequest struct {
	object.GenericCredentialOffer
	Form *response.RTUResponse `json:"form"`
}
