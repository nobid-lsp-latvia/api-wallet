// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	"crypto/sha256"
	"encoding/base64"

	"git.zzdats.lv/edim/api-wallet/attestation"
	"git.zzdats.lv/edim/api-wallet/routes/request"

	"azugo.io/azugo"
	"azugo.io/core/http"
	"github.com/valyala/fasthttp"
)

var b64 = base64.RawURLEncoding.Strict()

// @operationId WalletInstance
// @title Initialize wallet instance
// @description Initialize wallet instance using attestation key.
// @param AttestationRequest body request.AttestationRequest true "Attestation key request"
// @success 201 {empty} "Created"
// @failure 400 string string "Bad request"
// @failure 422 string string "Invalid request"
// @failure 500 string string "Internal server error"
// @resource Instance
// @route /instance [post].
func (r *router) attestation(ctx *azugo.Context) {
	req := request.AttestationRequest{}

	if err := ctx.Body.JSON(&req); err != nil {
		ctx.Error(err)

		return
	}

	buf, err := b64.DecodeString(req.Challenge)
	if err != nil {
		ctx.Error(azugo.ParamInvalidError{
			Name: "challenge",
			Tag:  "invalid",
			Err:  err,
		})

		return
	}

	attch := sha256.Sum256(buf)

	attest, err := r.Attestation().Verify(req.KeyAttestation, attch[:], req.HardwareKeyTag)
	if err != nil {
		ctx.Error(err)

		return
	}

	token, err := r.OpenID4VCI().ValidateNonce(ctx, req.Challenge)
	if err != nil {
		ctx.Error(err)

		return
	}

	if token == "anonymous" {
		if err := r.Store().Exec(ctx, "wallet.create_instance", &struct {
			*attestation.Result `json:",inline"`

			Person any `json:"person"`
		}{
			Result: attest,
			Person: struct {
				Code          string `json:"code"`
				RequesterCode string `json:"requesterCode"`
			}{
				Code:          token,
				RequesterCode: token,
			},
		}, nil); err != nil {
			ctx.Error(err)

			return
		}

		ctx.StatusCode(fasthttp.StatusCreated)

		return
	}

	session, err := r.IDAuth().UserInfo(ctx, http.WithHeader(fasthttp.HeaderAuthorization, "Bearer "+token))
	if err != nil {
		ctx.Error(err)

		return
	}

	if err := r.Store().Exec(ctx, "wallet.create_instance", &struct {
		*attestation.Result `json:",inline"`

		Person any `json:"person"`
	}{
		Result: attest,
		Person: struct {
			Code          string `json:"code"`
			GivenName     string `json:"givenName"`
			FamilyName    string `json:"familyName"`
			RequesterCode string `json:"requesterCode"`
		}{
			Code:          session.Code,
			GivenName:     session.GivenName,
			FamilyName:    session.FamilyName,
			RequesterCode: session.Code,
		},
	}, nil); err != nil {
		ctx.Error(err)

		return
	}

	ctx.StatusCode(fasthttp.StatusCreated)
}
