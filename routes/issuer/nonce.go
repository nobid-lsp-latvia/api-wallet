// SPDX-License-Identifier: EUPL-1.2

package issuer

import (
	"azugo.io/azugo"
)

func (r *router) nonce(ctx *azugo.Context) {
	nonce, err := r.OpenID4VCI().Nonce(ctx)
	if err != nil {
		ctx.Error(err)

		return
	}

	ctx.Header.Set("Cache-Control", "no-store")
	ctx.JSON(struct {
		Nonce string `json:"c_nonce"` //nolint: tagliatelle
	}{
		Nonce: nonce,
	})
}
