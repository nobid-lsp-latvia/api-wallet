// SPDX-License-Identifier: EUPL-1.2

package issuer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"azugo.io/azugo"
	"azugo.io/core/http"
	"github.com/valyala/fasthttp"
)

// @operationId GetOpenIDCredentialIssuer
// @title Gets OpenID Credential Issuer
// @description Gets openid-credential-issuer.json
// @success 200 object string "OK"
// @failure 400 string string "Bad request"
// @failure 401 {empty} "Unauthorized"
// @failure 403 {empty} "Forbidden"
// @failure 422 string string "Invalid request"
// @failure 500 string string "Internal server error"
// @resource WellKnown
// @route /.well-known/openid-credential-issuer [get].
func (r *router) openIDCredentialIssuer(ctx *azugo.Context) {
	client := ctx.HTTPClient().WithBaseURL(r.Config().Issuer.APIURL).WithOptions(&http.TLSConfig{
		InsecureSkipVerify: true,
	})

	var res map[string]any

	err := client.GetJSON("/.well-known/openid-credential-issuer", &res)
	if err != nil {
		ctx.Error(err)

		return
	}

	publicURL := strings.TrimSuffix(r.Config().WalletPublicURL, "/")

	res["nonce_endpoint"] = publicURL + "/nonce"
	res["credential_issuer"] = publicURL
	res["credential_endpoint"] = publicURL + "/credential"

	ctx.JSON(res)
}

// @operationId GetOpenIDConfiguration
// @title Gets OpenID Configuration
// @description Gets openid-configuration.json
// @success 200 object string "OK"
// @failure 400 string string "Bad request"
// @failure 401 {empty} "Unauthorized"
// @failure 403 {empty} "Forbidden"
// @failure 422 string string "Invalid request"
// @failure 500 string string "Internal server error"
// @resource WellKnown
// @route /.well-known/openid-configuration [get].
func (r *router) openIDConfiguration(ctx *azugo.Context) {
	client := ctx.HTTPClient().WithBaseURL(r.Config().Issuer.APIURL).WithOptions(&http.TLSConfig{
		InsecureSkipVerify: true,
	})

	var res map[string]any

	err := client.GetJSON("/.well-known/openid-configuration", &res)
	if err != nil {
		ctx.Error(err)

		return
	}

	publicURL := strings.TrimSuffix(r.Config().WalletPublicURL, "/")

	res["token_endpoint"] = publicURL + "/token"
	res["issuer"] = publicURL
	res["jwks_uri"] = publicURL + "/.well-known/jwks"

	ctx.JSON(res)
}

// @operationId GetOpenIDJWKS
// @title Gets OpenID JWKS
// @description Gets OpenID configuration JSON Web Key Set
// @success 200 object string "OK"
// @failure 400 string string "Bad request"
// @failure 401 {empty} "Unauthorized"
// @failure 403 {empty} "Forbidden"
// @failure 422 string string "Invalid request"
// @failure 500 string string "Internal server error"
// @resource WellKnown
// @route /.well-known/jwks [get].
func (r *router) openIDJWKS(ctx *azugo.Context) {
	client := ctx.HTTPClient().WithBaseURL(r.Config().Issuer.APIURL).WithOptions(&http.TLSConfig{
		InsecureSkipVerify: true,
	})

	res := struct {
		Keys []map[string]any `json:"keys"`
	}{}

	err := client.GetJSON("/static/jwks.json", &res)
	if err != nil {
		ctx.Error(err)

		return
	}

	kid, err := r.OpenID4VCI().SigningCertificateKID()
	if err != nil {
		ctx.Error(err)

		return
	}

	// Check if the kid is present in the keys
	var ok bool

	for _, key := range res.Keys {
		if key["kid"] == kid {
			ok = true

			break
		}
	}

	if !ok {
		pubKey, err := r.OpenID4VCI().SigningPublicKey()
		if err != nil {
			ctx.Error(err)

			return
		}

		key, err := publicKeyToJWK(pubKey)
		if err != nil {
			ctx.Error(err)

			return
		}

		key["kid"] = kid
		key["use"] = "sig"

		res.Keys = append(res.Keys, key)
	}

	ctx.JSON(res)
}

func (r *router) credential(ctx *azugo.Context) {
	client := ctx.HTTPClient()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	ctx.Request().CopyTo(req)

	req.SetRequestURI(r.Config().Issuer.APIURL + "/credential")
	req.Header.SetMethod("POST")

	resp := &ctx.Context().Response

	httpRequest := &http.Request{
		Request: req,
	}

	httpResponse := &http.Response{
		Response: resp,
	}

	err := client.Do(httpRequest, httpResponse)
	if err != nil {
		ctx.Error(err)

		return
	}

	ctx.Raw(httpResponse.Body())
	ctx.StatusCode(httpResponse.StatusCode())
}

func (r *router) token(ctx *azugo.Context) {
	grantType, err := ctx.Form.String("grant_type")
	if err != nil {
		ctx.Error(err)

		return
	}

	if grantType == "urn:ietf:params:oauth:grant-type:jwt-bearer" {
		if err := r.OpenID4VCI().Assertion(ctx); err != nil {
			d := azugo.BadRequestError{}
			if errors.As(err, &d) {
				ctx.StatusCode(fasthttp.StatusBadRequest)
				ctx.JSON(struct {
					Error string `json:"error"`
				}{
					Error: d.Description,
				})

				return
			}

			ctx.Error(err)

			return
		}

		return
	}

	client := ctx.HTTPClient().WithBaseURL(r.Config().Issuer.APIURL)

	data := make(map[string][]string)

	clientID, err := ctx.Form.String("client_id")
	if err != nil {
		ctx.Error(err)

		return
	}

	data["client_id"] = []string{clientID}

	data["grant_type"] = []string{grantType}

	preAuthorizedCode, err := ctx.Form.String("pre-authorized_code")
	if err != nil {
		ctx.Error(err)

		return
	}

	data["pre-authorized_code"] = []string{preAuthorizedCode}

	code, _ := ctx.Form.String("tx_code")
	if code == "" {
		code, err = r.App.Issuer().GetTXCode(ctx, preAuthorizedCode)
		if err != nil {
			ctx.Error(err)

			return
		}
	}

	data["tx_code"] = []string{code}

	res, err := client.PostForm("/token", data)
	if err != nil {
		ctx.Error(err)

		return
	}

	var jsonData any

	err = json.Unmarshal(res, &jsonData)
	if err != nil {
		ctx.Error(fmt.Errorf("error parsing JSON string: %w", err))

		return
	}

	ctx.JSON(jsonData)
}
