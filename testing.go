// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"testing"

	"github.com/go-quicktest/qt"
)

// TestApp for unit testing.
func TestApp(tb testing.TB) *App {
	tb.Helper()

	tb.Setenv("METRICS_ENABLED", "false")

	tb.Setenv("IDAUTH_URL", "http://idauth:8080")
	tb.Setenv("IDAUTH_CLIENT_ID", "edim.self-service.portal")
	tb.Setenv("IDAUTH_CLIENT_SECRET", "secret")

	tb.Setenv("ISSUER_NONCE_SHARED_SECRET", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")
	tb.Setenv("ISSUER_API_URL", "http://issuer:5000")
	tb.Setenv("MDL_API_URL", "http://mdl:5000")
	tb.Setenv("RTU_API_URL", "http://rtu:5000")
	tb.Setenv("FPRIS_API_URL", "http://fpris:5000")

	app, err := New(nil, "1.0.0-test")
	qt.Assert(tb, qt.IsNil(err))

	return app
}
