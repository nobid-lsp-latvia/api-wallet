// SPDX-License-Identifier: EUPL-1.2

//nolint:tagliatelle
package request

// AttestationRequest is a request model wallet initialization.
type AttestationRequest struct {
	// Challenge is a nonce that is used to verify the attestation.
	Challenge string `json:"challenge"`
	// KeyAttestation is a Base64 encoded CBOR attestation statement.
	KeyAttestation string `json:"key_attestation"`
	// HardwareKeyTag is a hardware key tag that is used to verify the attestation.
	HardwareKeyTag string `json:"hardware_key_tag"`
}
