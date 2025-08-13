// SPDX-License-Identifier: EUPL-1.2

package attestation

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"azugo.io/azugo"
	"github.com/fxamacker/cbor/v2"
)

// https://developer.android.com/privacy-and-security/security-key-attestation#root_certificate
const androidAppAttestRootCA = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAr7bHgiuxpwHsK7Qui8xU
FmOr75gvMsd/dTEDDJdSSxtf6An7xyqpRR90PL2abxM1dEqlXnf2tqw1Ne4Xwl5j
lRfdnJLmN0pTy/4lj4/7tv0Sk3iiKkypnEUtR6WfMgH0QZfKHM1+di+y9TFRtv6y
//0rb+T+W8a9nsNL/ggjnar86461qO0rOs2cXjp3kOG1FEJ5MVmFmBGtnrKpa73X
pXyTqRxB/M0n1n/W9nGqC4FSYa04T6N5RIZGBN2z2MT5IKGbFlbC8UrW0DxW7AYI
mQQcHtGl/m00QLVWutHQoVJYnFPlXTcHYvASLu+RhhsbDmxMgJJ0mcDpvsC4PjvB
+TxywElgS70vE0XmLD+OJtvsBslHZvPBKCOdT0MS+tgSOIfga+z1Z1g7+DVagf7q
uvmag8jfPioyKvxnK/EgsTUVi2ghzq8wm27ud/mIM7AY2qEORR8Go3TVB4HzWQgp
Zrt3i5MIlCaY504LzSRiigHCzAPlHws+W0rB5N+er5/2pJKnfBSDiCiFAVtCLOZ7
gLiMm0jhO2B6tUXHI/+MRPjy02i59lINMRRev56GKtcd9qO/0kUJWdZTdA2XoS82
ixPvZtXQpUpuL12ab+9EaDK8Z4RHJYYfCT3Q5vNAXaiWQ+8PTWm2QgBR/bkwSWc+
NpUFgNPN9PvQi8WEg5UmAGMCAwEAAQ==
-----END PUBLIC KEY-----`

func (a *Service) verifyAndroid(att string, challenge []byte, tag string, now time.Time) (*Result, error) {
	buf, err := base64.RawURLEncoding.DecodeString(att)
	if err != nil {
		return nil, azugo.ParamInvalidError{
			Name: "key_attestation",
			Tag:  "invalid",
			Err:  err,
		}
	}

	s := androidAttestation{}

	decoder := cbor.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(&s); err != nil {
		return nil, azugo.ParamInvalidError{
			Name: "key_attestation",
			Tag:  "invalid",
			Err:  err,
		}
	}

	cert, interms, err := a.verifyCert(s.AttStmt.X5c, androidAppAttestRootCA, now)
	if err != nil {
		return nil, azugo.ParamInvalidError{
			Name: "certificate",
			Tag:  "invalid",
			Err:  err,
		}
	}

	var isStrongbox bool

	for _, c := range interms {
		if cert.Issuer.String() != c.Subject.String() {
			continue
		}

		for _, name := range c.Subject.Names {
			val, ok := name.Value.(string)
			if !ok {
				continue
			}

			if val == "StrongBox" && (name.Type.Equal([]int{2, 5, 4, 10}) || name.Type.Equal([]int{2, 5, 4, 12})) {
				isStrongbox = true

				break
			}
		}

		if isStrongbox {
			break
		}
	}

	if !isStrongbox {
		return nil, azugo.ParamInvalidError{
			Name: "certificate",
			Tag:  "insecure",
		}
	}

	// Android Key Attestation provision information extansion data is identified by OID "1.3.6.1.4.1.11129.2.1.30"
	var provExtBytes []byte

	for _, ext := range cert.Extensions {
		if ext.Id.Equal([]int{1, 3, 6, 1, 4, 1, 11129, 2, 1, 30}) {
			provExtBytes = ext.Value
		}
	}

	var certsIssued int

	if len(provExtBytes) > 0 {
		decoder := cbor.NewDecoder(bytes.NewReader(provExtBytes))

		prov := provisioningInfo{}
		if err := decoder.Decode(&prov); err != nil {
			return nil, azugo.ParamInvalidError{
				Name: "certificate:extension",
				Tag:  "invalid",
				Err:  err,
			}
		}

		certsIssued = prov.CertsIssued
	}

	// Android Key Attestation attestation certificate's extension data is identified by the OID "1.3.6.1.4.1.11129.2.1.17"
	var attExtBytes []byte

	for _, ext := range cert.Extensions {
		if ext.Id.Equal([]int{1, 3, 6, 1, 4, 1, 11129, 2, 1, 17}) {
			attExtBytes = ext.Value
		}
	}

	if len(attExtBytes) == 0 {
		return nil, azugo.ParamInvalidError{
			Name: "certificate:extension",
			Tag:  "not_found",
		}
	}

	decoded := keyDescription{}

	_, err = asn1.Unmarshal(attExtBytes, &decoded)
	if err != nil {
		return nil, azugo.ParamInvalidError{
			Name: "certificate:extension",
			Tag:  "invalid",
			Err:  err,
		}
	}

	if !bytes.Equal(decoded.AttestationChallenge, challenge) {
		return nil, azugo.ParamInvalidError{
			Name: "challenge",
			Tag:  "invalid",
		}
	}

	// Convert public key to X9.62 format
	pk, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type: %T", cert.PublicKey)
	}

	pubkey, err := pk.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key: %w", err)
	}

	calulatedTagBytes := sha256.Sum256(pubkey.Bytes())

	tagBytes, err := base64.StdEncoding.DecodeString(tag)
	if err != nil {
		return nil, azugo.ParamInvalidError{
			Name: "hardware_key_tag",
			Tag:  "invalid",
			Err:  err,
		}
	}

	// Verify the tag
	if !bytes.Equal(tagBytes, calulatedTagBytes[:]) {
		return nil, azugo.ParamInvalidError{
			Name: "hardware_key_tag",
			Tag:  "invalid",
		}
	}

	// Export the public key in PEM format
	buf, err = x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal public key: %w", err)
	}

	publicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: buf,
	})

	return &Result{
		DeviceType:     "android",
		CertsIssued:    certsIssued,
		HardwareKeyTag: tag,
		PublicKey:      string(publicKey),
	}, nil
}
