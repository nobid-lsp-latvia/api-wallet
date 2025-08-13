// SPDX-License-Identifier: EUPL-1.2

package openid4vci

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"math/big"
)

func (s *Service) publicKeyFromJWK(cnf map[string]any) (any, error) {
	jwk, ok := cnf["jwk"].(map[string]any)
	if !ok {
		return nil, errors.New("unsupported CNF")
	}

	kt, ok := jwk["kty"].(string)
	if !ok {
		return nil, errors.New("missing key type")
	}

	switch kt {
	case "EC":
		crv, ok := jwk["crv"].(string)
		if !ok {
			return nil, errors.New("missing curve")
		}

		var x, y *big.Int

		var xb string

		if xb, ok = jwk["x"].(string); !ok {
			return nil, errors.New("missing x")
		}

		xbuf, err := base64.RawURLEncoding.DecodeString(xb)
		if err != nil {
			return nil, err
		}

		x = new(big.Int).SetBytes(xbuf)

		var yb string

		if yb, ok = jwk["y"].(string); !ok {
			return nil, errors.New("missing y")
		}

		ybuf, err := base64.RawURLEncoding.DecodeString(yb)
		if err != nil {
			return nil, err
		}

		y = new(big.Int).SetBytes(ybuf)

		switch crv {
		case "P-256":
			return &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     x,
				Y:     y,
			}, nil
		case "P-384":
			return &ecdsa.PublicKey{
				Curve: elliptic.P384(),
				X:     x,
				Y:     y,
			}, nil
		case "P-521":
			return &ecdsa.PublicKey{
				Curve: elliptic.P521(),
				X:     x,
				Y:     y,
			}, nil
		default:
			return nil, errors.New("unsupported curve")
		}
	case "RSA":
		var (
			e int
			n *big.Int
		)

		var eb string

		if eb, ok = jwk["e"].(string); !ok {
			return nil, errors.New("missing exponent")
		}

		ebuf, err := base64.RawURLEncoding.DecodeString(eb)
		if err != nil {
			return nil, err
		}

		e = int(new(big.Int).SetBytes(ebuf).Uint64()) //nolint:gosec

		var nb string

		if nb, ok = jwk["n"].(string); !ok {
			return nil, errors.New("missing modulus")
		}

		nbuf, err := base64.RawURLEncoding.DecodeString(nb)
		if err != nil {
			return nil, err
		}

		n = new(big.Int).SetBytes(nbuf)

		return rsa.PublicKey{
			E: e,
			N: n,
		}, nil
	default:
		return nil, errors.New("unsupported key type")
	}
}
