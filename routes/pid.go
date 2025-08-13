// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	"errors"
	"fmt"

	"git.zzdats.lv/edim/api-wallet/models"
	"git.zzdats.lv/edim/api-wallet/routes/object"
	"git.zzdats.lv/edim/api-wallet/routes/request"

	"azugo.io/azugo"
	"azugo.io/core/http"
	"github.com/valyala/fasthttp"
)

func (r *router) qrCodeInternal(ctx *azugo.Context) {
	requestType := ctx.Params.String("requestType")

	r.qrCodeData(ctx, requestType, true)
}

// @operationId GenerateQRCode
// @title Generate QR code
// @description Generates qr code
// @param requestType path string true "options: `pid`, `rtu`, `mdl`"
// @success 200 GenerateCredentialOfferResponse models.GenerateCredentialOfferResponse "pid result"
// @failure 400 string string "Bad request"
// @failure 401 {empty} "Unauthorized"
// @failure 403 {empty} "Forbidden"
// @failure 422 string string "Invalid request"
// @failure 500 string string "Internal server error"
// @resource QRCode
// @route /1.0/{requestType} [post].
func (r *router) qrCode(ctx *azugo.Context) {
	requestType := ctx.Params.String("requestType")

	r.qrCodeData(ctx, requestType, false)
}

func (r *router) qrCodeData(ctx *azugo.Context, requestType string, showTXCode bool) {
	personCodeClaim := ctx.User().Claim("code")

	var (
		gcoReq any

		err error
	)

	if len(personCodeClaim) == 0 || personCodeClaim[0] == "" {
		ctx.StatusCode(fasthttp.StatusUnauthorized)

		return
	}

	switch requestType {
	case "pid":
		gcoReq, err = r.getPIDData(ctx)
		if err != nil {
			if errors.Is(err, http.NotFoundError{}) {
				ctx.Error(err)
				ctx.Text("Data about person not found")

				return
			}

			ctx.Error(fmt.Errorf("failed to call PID: %w", err))

			return
		}
	case "rtu":
		gcoReq, err = r.getRTUData(ctx)
		if err != nil {
			if errors.Is(err, http.NotFoundError{}) {
				ctx.Error(err)
				ctx.Text("Data about education not found")

				return
			}

			ctx.Error(fmt.Errorf("failed to call RTU: %w", err))

			return
		}
	case "mdl":
		gcoReq, err = r.getMDLData(ctx)
		if err != nil {
			if errors.Is(err, http.NotFoundError{}) {
				ctx.Error(err)
				ctx.Text("Data about drivers licence not found")

				return
			}

			ctx.Error(fmt.Errorf("failed to call MDL: %w", err))

			return
		}
	default:
		ctx.StatusCode(fasthttp.StatusNotFound)

		return
	}

	issuerRes := &models.GenerateCredentialOffer{}
	client := ctx.HTTPClient().WithBaseURL(r.Config().Issuer.APIURL)

	err = client.PostJSON("/generate_credential_offer", &gcoReq, issuerRes)
	if err != nil {
		ctx.Error(err)

		return
	}

	// ignore lint here, we need this to have reference in swagger
	res := &models.GenerateCredentialOffer{} //nolint:ineffassign,wastedassign

	res, err = r.App.Issuer().ParseCredentialOffer(ctx, *issuerRes, showTXCode)
	if err != nil {
		ctx.Error(err)

		return
	}

	ctx.JSON(res)
}

func (r *router) getPIDData(ctx *azugo.Context) (any, error) {
	// call fpris and get person data
	fprisRes, err := r.App.FprisAPIClient().GetFprisData(ctx)
	if err != nil {
		return nil, err
	}

	if fprisRes == nil {
		return nil, http.NotFoundError{Resource: "person"}
	}

	gcoReq := &request.PIDCredentialOfferRequest{
		GenericCredentialOffer: object.GenericCredentialOffer{
			CredentialIDS:      []string{"eu.europa.ec.eudi.pid_mdoc"},
			CodeGrant:          "pre_auth_code",
			CredentialOfferURI: r.Config().QRAPIDeepLink + "://",
			ReturnInHTML:       false,
		},
		Form: &request.PIDCredentialOfferForm{},
	}

	gcoReq.Form.FamilyName = fprisRes.FamilyName
	gcoReq.Form.GivenName = fprisRes.GivenName
	gcoReq.Form.BirthDate = fprisRes.BirthDate
	gcoReq.Form.PersonalAdministrativeNumber = fprisRes.PersonalAdministrativeNumber
	gcoReq.Form.BirthPlace = fprisRes.BirthPlace
	gcoReq.Form.BirthCountry = fprisRes.BirthCountry
	gcoReq.Form.EstimatedIssuanceDate = fprisRes.IssuanceDate
	gcoReq.Form.EstimatedExpiryDate = fprisRes.ExpiryDate
	gcoReq.Form.IssuingCountry = fprisRes.IssuingCountry
	gcoReq.Form.IssuingAuthority = fprisRes.IssuingAuthority
	gcoReq.Form.AgeOver18 = fprisRes.AgeOver18
	gcoReq.Form.Nationality = fprisRes.Nationality

	return gcoReq, nil
}

func (r *router) getRTUData(ctx *azugo.Context) (any, error) {
	// call rtu and get diploma data
	rtuRes, err := r.App.RTUAPIClient().GetRTUData(ctx)
	if err != nil {
		return nil, err
	}

	if rtuRes == nil {
		return nil, http.NotFoundError{Resource: "person"}
	}

	gcoReq := &request.RTUCredentialOfferRequest{
		GenericCredentialOffer: object.GenericCredentialOffer{
			CredentialIDS:      []string{"eu.europa.ec.eudi.rtu_diploma_mdoc"},
			CodeGrant:          "pre_auth_code",
			CredentialOfferURI: r.Config().QRAPIDeepLink + "://",
			ReturnInHTML:       false,
		},
		Form: rtuRes,
	}

	return gcoReq, nil
}

func (r *router) getMDLData(ctx *azugo.Context) (any, error) {
	// call mdl and get driving license data
	mdlRes, err := r.App.MDLAPIClient().GetMDLData(ctx)
	if err != nil {
		return nil, err
	}

	if mdlRes == nil {
		return nil, http.NotFoundError{Resource: "person"}
	}

	gcoReq := &request.MDLCredentialOfferRequest{
		GenericCredentialOffer: object.GenericCredentialOffer{
			CredentialIDS:      []string{"eu.europa.ec.eudi.mdl_mdoc"},
			CodeGrant:          "pre_auth_code",
			CredentialOfferURI: r.Config().QRAPIDeepLink + "://",
			ReturnInHTML:       false,
		},
		Form: mdlRes,
	}

	return gcoReq, nil
}
