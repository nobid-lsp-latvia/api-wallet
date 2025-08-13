// SPDX-License-Identifier: EUPL-1.2

package attestation

import (
	"crypto/subtle"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"
)

func (a *Service) verifyCert(chain [][]byte, rootCA string, now time.Time) (*x509.Certificate, []*x509.Certificate, error) {
	roots := x509.NewCertPool()

	var pub any

	if strings.HasPrefix(rootCA, "-----BEGIN PUBLIC KEY-----") {
		var err error

		block, _ := pem.Decode([]byte(rootCA))

		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, nil, err
		}
	} else {
		roots.AppendCertsFromPEM([]byte(rootCA))
	}

	interms := make([]*x509.Certificate, 0, len(chain)-2)

	var cert *x509.Certificate

	for _, buf := range chain {
		c, err := x509.ParseCertificate(buf)
		if err != nil {
			return nil, nil, err
		}

		// TODO check CRLs

		if c.Subject.String() == c.Issuer.String() {
			if pub != nil && equalKeys(c.PublicKey, pub) {
				roots.AddCert(c)
			}

			continue
		}

		// KeyUsageDigitalSignature is used by both Android and iOS
		if c.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
			cert = c

			continue
		}

		interms = append(interms, c)
	}

	intermediates := x509.NewCertPool()

	for _, c := range interms {
		intermediates.AddCert(c)
	}

	if cert == nil {
		return nil, nil, errors.New("valid certificate not found in attestation")
	}

	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		CurrentTime:   now,
	}); err != nil {
		return nil, nil, err
	}

	return cert, interms, nil
}

func equalKeys(a, b interface{}) bool {
	aBytes, err := x509.MarshalPKIXPublicKey(a)
	if err != nil {
		return false
	}

	bBytes, err := x509.MarshalPKIXPublicKey(b)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(aBytes, bBytes) == 1
}
