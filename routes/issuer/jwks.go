// SPDX-License-Identifier: EUPL-1.2

package issuer

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"math/big"
)

func publicKeyToJWK(publicKey any) (map[string]any, error) {
	switch pk := publicKey.(type) {
	case *ecdsa.PublicKey:
		x := pk.X.Bytes()
		y := pk.Y.Bytes()

		if len(x) > 32 {
			x = x[len(x)-32:]
		}

		if len(y) > 32 {
			y = y[len(y)-32:]
		}

		return map[string]any{
			"kty": "EC",
			"crv": pk.Curve.Params().Name,
			"x":   base64.RawURLEncoding.EncodeToString(x),
			"y":   base64.RawURLEncoding.EncodeToString(y),
		}, nil
	case *rsa.PublicKey:
		n := pk.N.Bytes()
		e := big.NewInt(int64(pk.E)).Bytes()

		if len(n) > 256 {
			n = n[len(n)-256:]
		}

		if len(e) > 3 {
			e = e[len(e)-3:]
		}

		return map[string]any{
			"kty": "RSA",
			"n":   base64.RawURLEncoding.EncodeToString(n),
			"e":   base64.RawURLEncoding.EncodeToString(e),
		}, nil
	default:
		return nil, errors.New("unsupported key type")
	}
}
