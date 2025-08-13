// SPDX-License-Identifier: EUPL-1.2

//nolint:tagliatelle
package response

import "time"

// RTUResponse contains the data of the diploma.
type RTUResponse struct {
	// GivenName represents the person given name
	GivenName string `json:"givenName"`
	// FamilyName represents the person family name
	FamilyName string `json:"familyName"`
	// NationalID represents the person national ID
	NationalID string `json:"nationalID"`
	// CitizenshipCountryCode represents the citizenship country code
	CitizenshipCountryCode string `json:"citizenshipCountryCode"`
	// Type represents the type of the diploma
	Type string `json:"type"`
	// Title represents the title
	Title string `json:"title"`
	// AwardingDate  represents the awarding date of the diploma
	AwardingDate time.Time `json:"awardingDate"`
	// AwardingBodyLegalName represents the awarding body legal name
	AwardingBodyLegalName string `json:"awardingBody_legalName"`
	// AwardingBodyRegistration represents the awarding body registration
	AwardingBodyRegistration string `json:"awardingBody_registration"`
	// AwardingBodyCountryCode represents the awarding body country code
	AwardingBodyCountryCode string `json:"awardingBody_countryCode"`
	// CreditPoint represents the credit point
	CreditPoint string `json:"creditPoint"`
	// CreditValue represents the credit value
	CreditValue string `json:"creditValue"`
	// MaximumDuration represents the maximum duration
	MaximumDuration string `json:"maximumDuration"`
	// ThematicArea represents the thematic area
	ThematicArea string `json:"thematicArea"`
	// NqfLevel reprsents the nqf level
	NqfLevel string `json:"nqfLevel"`
	// EqfLevel reprsents the nqf level
	EqfLevel string `json:"eqfLevel"`
	// ValidFrom reprsents the date the diploma is valid from
	ValidFrom time.Time `json:"validFrom"`
	// IssuanceDate reprsents the date of the diploma issuance
	IssuanceDate time.Time `json:"issuanceDate"`
	// Issued reprsents the date the diploma was issued
	Issued time.Time `json:"issued"`
}
