// SPDX-License-Identifier: EUPL-1.2

//nolint:tagliatelle
package response

import "git.zzdats.lv/edim/api-wallet/util"

// CategoryRestriction defines the driving privilege category restriction.
type CategoryRestriction struct {
	// Sign as per ISO/IEC 18013-2 Annex A
	Sign string `json:"sign"`
	// Value as per ISO/IEC 18013-2 Annex A
	Value string `json:"value"`
}

// DrivingPrivilege defines driving license categories.
type DrivingPrivilege struct {
	// Driver category code
	VehicleCategoryCode string `json:"vehicle_category_code"`
	// Starting date of validity of the driver category
	IssueDate util.Date `json:"issue_date"`
	// Driver category expiry date
	ExpiryDate util.Date `json:"expiry_date"`
	// Category restrictions
	Code []CategoryRestriction `json:"code"`
}

// MDLResponse defines the response structure for the CSDD data.
type MDLResponse struct {
	// PersonalAdministrativeNumber represents driver's administrative number
	PersonalAdministrativeNumber string `json:"personal_administrative_number"`
	// DocumentNumber represents document certificate number
	DocumentNumber string `json:"document_number"`
	// BirthData represents driver's date of birth
	BirthDate util.Date `json:"birth_date"`
	// GivenName represents driver's name
	GivenName string `json:"given_name"`
	// FamilyName represents driver's surname
	FamilyName string `json:"family_name"`
	// IssueDate represents driver's license start date of validity
	IssueDate util.Date `json:"issue_date"`
	// ExpireDate reprsents driver's license expiration date
	ExpiryDate util.Date `json:"expiry_date"`
	// IssuingCountry represents driving license issuing country code (ISO 3166-1 alpha-2)
	IssuingCountry string `json:"issuing_country"`
	// IssuingAuthority represents issuing authority
	IssuingAuthority string `json:"issuing_authority"`
	// DrivingPrivileges represents driving license categories
	DrivingPrivileges []DrivingPrivilege `json:"driving_privileges"`
	// UnDistinguishingSign represents distinguishing sign of the issuing country according to ISO/IEC 18013-1:2018
	UnDistinguishingSign string `json:"un_distinguishing_sign"`
	// Portrait represent photo of the driver of the vehicle
	Portrait string `json:"portrait"`
}
