// SPDX-License-Identifier: EUPL-1.2

package openid4vci

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"azugo.io/azugo"
	"azugo.io/core/http"
	"github.com/golang-jwt/jwt/v5"
	jsondb "github.com/nobid-lsp-latvia/lx-go-jsondb"
	"go.uber.org/zap"
)

type AttestationPerson struct {
	Code       string `json:"code"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

type WalletAttestation struct {
	Format            string `json:"format"`
	WalletAttestation string `json:"wallet_attestation"`
}

func (s *Service) Assertion(ctx *azugo.Context) error {
	assertion, err := ctx.Form.String("assertion")
	if err != nil {
		return err
	}

	var instanceID string
	// var person *AttestationPerson

	token, err := jwt.Parse(assertion, func(t *jwt.Token) (any, error) {
		// TODO: use kid or hardware_key_tag or empheral key pub?
		keyTag, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing key ID")
		}

		resp := struct {
			PublicKey string             `json:"publicKey"`
			Person    *AttestationPerson `json:"person"`
		}{}

		if err := s.store.Exec(ctx, "wallet.get_public_key", &struct {
			Type           string `json:"type"`
			HardwareKeyTag string `json:"hardwareKeyTag"`
		}{
			Type:           "instance",
			HardwareKeyTag: keyTag,
		}, &resp); err != nil {
			var eerr jsondb.ExecError
			if errors.As(err, &eerr) && eerr.Code == "err:public_key:not_found" {
				return nil, errors.New("public key not found")
			}

			s.app.Log().Error("failed to get public key", zap.Error(err))

			return nil, errors.New("failed to get public key")
		}

		instanceID, err = url.JoinPath(s.walletInstanceURL, keyTag)
		if err != nil {
			return nil, fmt.Errorf("failed to generate instance ID: %w", err)
		}
		// person = resp.Person

		// Parse PEM encoded public key
		block, _ := pem.Decode([]byte(resp.PublicKey))
		if block == nil {
			return nil, errors.New("failed to parse public key")
		}

		return x509.ParsePKIXPublicKey(block.Bytes)
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}),
		// TODO: add leeway config
		// jwt.WithLeeway(s.config.)
		jwt.WithExpirationRequired(),
		jwt.WithSubject(s.walletPublicURL),
	)
	if err != nil {
		if isJWTError(err) {
			return azugo.BadRequestError{Description: err.Error()}
		}

		return err
	}

	if !token.Valid {
		return azugo.BadRequestError{Description: "invalid token"}
	}

	if typ, ok := token.Header["typ"].(string); !ok || typ != "var+jwt" {
		return azugo.BadRequestError{Description: "invalid token type for wallet attestation"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return azugo.BadRequestError{Description: "invalid claims"}
	}

	if typ, ok := claims["type"].(string); !ok || typ != "WalletInstanceAttestationRequest" {
		return azugo.BadRequestError{Description: "invalid token type for wallet attestation"}
	}

	// Check if valid issuer
	issuer, err := claims.GetIssuer()
	if err != nil || issuer != instanceID {
		return azugo.BadRequestError{Description: "unknown issuer"}
	}

	cnf, ok := claims["cnf"].(map[string]any)
	if !ok {
		return azugo.BadRequestError{Description: "invalid or missing cnf"}
	}

	publicKey, err := s.publicKeyFromJWK(cnf)
	if err != nil {
		return azugo.BadRequestError{Description: "invalid or missing cnf", Err: err}
	}

	// Convert public key to X9.62 format
	pk, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return azugo.BadRequestError{Description: fmt.Sprintf("invalid public key type: %T", publicKey)}
	}

	pubkey, err := pk.ECDH()
	if err != nil {
		return fmt.Errorf("failed to convert public key: %w", err)
	}

	calulatedTagBytes := sha256.Sum256(pubkey.Bytes())

	tag, ok := claims["hardware_key_tag"].(string)
	if !ok {
		return azugo.ParamInvalidError{
			Name: "hardware_key_tag",
			Tag:  "missing",
		}
	}

	tagBytes, err := base64.StdEncoding.DecodeString(tag)
	if err != nil {
		return azugo.ParamInvalidError{
			Name: "hardware_key_tag",
			Tag:  "invalid",
			Err:  err,
		}
	}

	// Verify the tag
	if !bytes.Equal(tagBytes, calulatedTagBytes[:]) {
		return azugo.ParamInvalidError{
			Name: "hardware_key_tag",
			Tag:  "invalid",
		}
	}

	// TODO: validate Android/Apple assertion
	// `nonce`, `hardware_signature` = sign(sha256(nonce+keypub)+hardware_key_tag)) and `key_attestation` claims

	// Currently do not include person data in the attestation
	return s.IssueAttestations(ctx, token, instanceID, nil /*person*/)
}

func (s *Service) IssueAttestations(ctx *azugo.Context, req *jwt.Token, instanceID string, person *AttestationPerson) error {
	tokc, _ := req.Claims.(jwt.MapClaims)

	token := jwt.New(jwt.SigningMethodES256)

	kid, err := s.SigningCertificateKID()
	if err != nil {
		return err
	}

	token.Header["typ"] = "oauth-client-attestation+jwt"
	token.Header["kid"] = kid

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"iss":         s.walletPublicURL,
		"sub":         s.walletPublicURL,
		"instance_id": instanceID,
		"cnf":         tokc["cnf"],
		"iat":         now.Unix(),
		"exp":         now.Add(24 * time.Hour).Unix(),
	}

	if person != nil {
		claims["personal_administrative_number"] = person.Code
		claims["given_name"] = person.GivenName
		claims["family_name"] = person.FamilyName
	}

	token.Claims = claims

	cert, err := s.SigningCertificate()
	if err != nil {
		return err
	}

	tok, err := token.SignedString(cert.PrivateKey)
	if err != nil {
		return err
	}

	// TODO: add dc+sd-jwt and mso_mdoc attestation tokens

	ctx.JSON(struct {
		WalletAttestations []WalletAttestation `json:"wallet_attestations"`
	}{
		WalletAttestations: []WalletAttestation{
			{
				Format:            "jwt",
				WalletAttestation: tok,
			},
			// {
			// 	Format:            "dc+sd-jwt",
			// 	WalletAttestation: "",
			// },
			// {
			// 	Format:            "mso_mdoc",
			// 	WalletAttestation: "",
			// },
		},
	})

	return nil
}

func (s *Service) VerifyAttestation(ctx *azugo.Context, tok string) (string, *AttestationPerson, error) {
	token, err := jwt.Parse(tok, func(_ *jwt.Token) (any, error) {
		return s.SigningPublicKey()
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}),
		// TODO: add leeway config
		// jwt.WithLeeway(s.config.)
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(s.walletPublicURL),
		jwt.WithSubject(s.walletPublicURL),
	)
	if err != nil {
		if isJWTError(err) {
			return "", nil, azugo.BadRequestError{Description: err.Error()}
		}

		return "", nil, err
	}

	if !token.Valid {
		return "", nil, azugo.BadRequestError{Description: "invalid token"}
	}

	if typ, ok := token.Header["typ"].(string); !ok || typ != "oauth-client-attestation+jwt" {
		return "", nil, azugo.BadRequestError{Description: "invalid token type"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, azugo.BadRequestError{Description: "invalid claims"}
	}

	sub, ok := claims["instance_id"].(string)
	if !ok {
		return "", nil, azugo.BadRequestError{Description: "invalid subject"}
	}

	keyTag := strings.TrimPrefix(strings.TrimPrefix(sub, s.walletInstanceURL), "/")
	if keyTag == "" {
		return "", nil, azugo.BadRequestError{Description: "invalid subject"}
	}

	resp := struct {
		ID             string             `json:"id"`
		Status         string             `json:"status"`
		HardwareKeyTag string             `json:"hardwareKeyTag"`
		FID            string             `json:"firebaseId"`
		Person         *AttestationPerson `json:"person"`
	}{}

	if err := s.store.Exec(ctx, "wallet.get_instance_by_tag", &struct {
		Type           string `json:"type"`
		HardwareKeyTag string `json:"hardwareKeyTag"`
	}{
		Type:           "instance",
		HardwareKeyTag: keyTag,
	}, &resp); err != nil {
		var eerr jsondb.ExecError
		if errors.As(err, &eerr) && eerr.Code == "err:instance:not_found" {
			return "", nil, http.NotFoundError{Resource: "instance"}
		}

		return "", nil, fmt.Errorf("failed to get wallet instance: %w", err)
	}

	return resp.ID, resp.Person, nil
}
