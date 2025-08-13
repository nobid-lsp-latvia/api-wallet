// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	"testing"

	api "git.zzdats.lv/edim/api-wallet"

	"azugo.io/azugo"
	"github.com/go-quicktest/qt"
)

func testApp(t testing.TB) *azugo.TestApp {
	app := api.TestApp(t)

	err := Init(app)
	qt.Assert(t, qt.IsNil(err))

	return azugo.NewTestApp(app.App)
}
