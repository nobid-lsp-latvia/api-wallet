// SPDX-License-Identifier: EUPL-1.2

package attestation

import (
	"encoding/asn1"
)

type keyDescription struct {
	AttestationVersion       int
	AttestationSecurityLevel asn1.Enumerated
	KeymasterVersion         int
	KeymasterSecurityLevel   asn1.Enumerated
	AttestationChallenge     []byte
	UniqueID                 []byte
	SoftwareEnforced         authorizationList
	TeeEnforced              authorizationList
}

type authorizationList struct {
	Purpose                     []int       `asn1:"tag:1,explicit,set,optional"`
	Algorithm                   int         `asn1:"tag:2,explicit,optional"`
	KeySize                     int         `asn1:"tag:3,explicit,optional"`
	Digest                      []int       `asn1:"tag:5,explicit,set,optional"`
	Padding                     []int       `asn1:"tag:6,explicit,set,optional"`
	EcCurve                     int         `asn1:"tag:10,explicit,optional"`
	RsaPublicExponent           int         `asn1:"tag:200,explicit,optional"`
	MgfDigest                   []int       `asn1:"tag:203,explicit,set,optional"`
	RollbackResistance          any         `asn1:"tag:303,explicit,optional"`
	EarlyBootOnly               any         `asn1:"tag:305,explicit,optional"`
	ActiveDateTime              int         `asn1:"tag:400,explicit,optional"`
	OriginationExpireDateTime   int         `asn1:"tag:401,explicit,optional"`
	UsageExpireDateTime         int         `asn1:"tag:402,explicit,optional"`
	UsageCountLimit             int         `asn1:"tag:403,explicit,optional"`
	NoAuthRequired              any         `asn1:"tag:503,explicit,optional"`
	UserAuthType                int         `asn1:"tag:504,explicit,optional"`
	AuthTimeout                 int         `asn1:"tag:505,explicit,optional"`
	AllowWhileOnBody            any         `asn1:"tag:506,explicit,optional"`
	TrustedUserPresenceRequired any         `asn1:"tag:507,explicit,optional"`
	TrustedConfirmationRequired any         `asn1:"tag:508,explicit,optional"`
	UnlockedDeviceRequired      any         `asn1:"tag:509,explicit,optional"`
	AllApplications             any         `asn1:"tag:600,explicit,optional"`
	ApplicationID               any         `asn1:"tag:601,explicit,optional"`
	CreationDateTime            int         `asn1:"tag:701,explicit,optional"`
	Origin                      int         `asn1:"tag:702,explicit,optional"`
	RootOfTrust                 rootOfTrust `asn1:"tag:704,explicit,optional"`
	OsVersion                   int         `asn1:"tag:705,explicit,optional"`
	OsPatchLevel                int         `asn1:"tag:706,explicit,optional"`
	AttestationApplicationID    []byte      `asn1:"tag:709,explicit,optional"`
	AttestationIDBrand          []byte      `asn1:"tag:710,explicit,optional"`
	AttestationIDDevice         []byte      `asn1:"tag:711,explicit,optional"`
	AttestationIDProduct        []byte      `asn1:"tag:712,explicit,optional"`
	AttestationIDSerial         []byte      `asn1:"tag:713,explicit,optional"`
	AttestationIDImei           []byte      `asn1:"tag:714,explicit,optional"`
	AttestationIDMeid           []byte      `asn1:"tag:715,explicit,optional"`
	AttestationIDManufacturer   []byte      `asn1:"tag:716,explicit,optional"`
	AttestationIDModel          []byte      `asn1:"tag:717,explicit,optional"`
	VendorPatchLevel            int         `asn1:"tag:718,explicit,optional"`
	BootPatchLevel              int         `asn1:"tag:719,explicit,optional"`
	DeviceUniqueAttestation     any         `asn1:"tag:720,explicit,optional"`
	AttestationIDSecondIMEI     []byte      `asn1:"tag:723,explicit,optional"`
	ModuleHash                  []byte      `asn1:"tag:724,explicit,optional"`
}

type rootOfTrust struct {
	VerifiedBootKey   []byte
	DeviceLocked      bool
	VerifiedBootState verifiedBootState
	VerifiedBootHash  []byte
}

type verifiedBootState int

const (
	Verified verifiedBootState = iota
	SelfSigned
	Unverified
	Failed
)

type provisioningInfo struct {
	CertsIssued int `cbor:"1"` // '1' corresponds to the OID field
}

type androidAttestation struct {
	AttStmt *androidStmt `cbor:"attStmt"`
}

type androidStmt struct {
	X5c [][]byte `cbor:"x5c"`
}
