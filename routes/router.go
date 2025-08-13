// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	"errors"
	"strings"

	wallet "git.zzdats.lv/edim/api-wallet"
	"git.zzdats.lv/edim/api-wallet/openapi"
	"git.zzdats.lv/edim/api-wallet/routes/issuer"

	"azugo.io/azugo"
	"azugo.io/azugo/token"
	"azugo.io/azugo/user"
	"azugo.io/core/http"
	"github.com/nobid-lsp-latvia/go-idauth"
	oa "github.com/nobid-lsp-latvia/go-openapi"
	"github.com/valyala/fasthttp"
)

type router struct {
	*wallet.App
	openapi *oa.OpenAPI
}

func Init(a *wallet.App) error {
	r := &router{
		App: a,
	}
	r.openapi = oa.NewDefaultOpenAPIHandler(openapi.OpenAPIDefinition, a.App)

	a.Get("/healthz", r.healthz)

	if err := issuer.Bind(a, a); err != nil {
		return err
	}

	a.Post("/instance", r.attestation)

	a.Post("/eparaksts/{type}", r.eparakstsPrepare)
	a.Get("/eparaksts/download/{requestID}", r.eparakstsDownload)
	a.Delete("/eparaksts/{type}/{requestID}", r.eparakstsCloseSession)
	a.Post("/eparaksts/validate", r.eparakstsValidate)

	a.Get("/eparaksts/identities", r.eparakstsIdentities)
	a.Get("/eparaksts/identities/{id}", r.eparakstsIdentitiesUlid)

	portal := a.Group("/1.0/internal")
	{
		portal.Use(idauth.Authentication(a.App, a.Config().IDAuth))
		portal.Post("/{requestType}", idauth.UserHasScope("citizen", r.qrCodeInternal))
	}

	v1 := a.Group("/1.0")
	{
		v1.Use(idauth.Authentication(a.App, a.Config().IDAuth))
		v1.Post("/{requestType}", r.qrCode)
	}

	return nil
}

// todo: might use for future
//
//nolint:unused
func (r *router) authentication(next azugo.RequestHandler) azugo.RequestHandler {
	return func(ctx *azugo.Context) {
		auth := ctx.Header.Get(fasthttp.HeaderAuthorization)
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			ctx.Error(http.UnauthorizedError{})

			return
		}

		tok := strings.TrimPrefix(auth, "Bearer ")

		instanceID, person, err := r.App.OpenID4VCI().VerifyAttestation(ctx, tok)
		if err != nil {
			if errors.Is(err, azugo.BadRequestError{}) || errors.Is(err, http.NotFoundError{}) {
				ctx.Error(http.UnauthorizedError{})

				return
			}

			ctx.Error(err)

			return
		}

		ctx.SetUser(user.New(map[string]token.ClaimStrings{
			"sub":         {instanceID},
			"code":        {person.Code},
			"given_name":  {person.GivenName},
			"family_name": {person.FamilyName},
			"scope":       {"citizen"},
		}))

		next(ctx)
	}
}
