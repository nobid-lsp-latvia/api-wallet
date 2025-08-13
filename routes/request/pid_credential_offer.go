// SPDX-License-Identifier: EUPL-1.2

//nolint:tagliatelle
package request

import "git.zzdats.lv/edim/api-wallet/routes/object"

type PIDCredentialOfferRequest struct {
	object.GenericCredentialOffer
	Form *PIDCredentialOfferForm `json:"form"`
}

type PIDCredentialOfferForm struct {
	object.FprisData
	EstimatedIssuanceDate string `json:"estimated_issuance_date"`
	EstimatedExpiryDate   string `json:"estimated_expiry_date"`
}
