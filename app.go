// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"git.zzdats.lv/edim/api-wallet/attestation"
	"git.zzdats.lv/edim/api-wallet/issuer"
	"git.zzdats.lv/edim/api-wallet/openid4vci"
	"git.zzdats.lv/edim/api-wallet/tasks"
	jsondb "github.com/nobid-lsp-latvia/lx-go-jsondb"

	"azugo.io/azugo"
	"azugo.io/azugo/server"
	"github.com/nobid-lsp-latvia/go-idauth"
	"github.com/spf13/cobra"
)

// App is the application instance.
type App struct {
	*azugo.App

	config *Configuration

	store jsondb.Store

	vci         *openid4vci.Service
	issuer      *issuer.Issuer
	attestation *attestation.Service
	idauth      *idauth.Client

	fprisAPIClient   *FprisAPIClient
	rtuAPIClient     *RTUAPIClient
	mdlAPIClient     *MDLAPIClient
	simpleSignClient *SimpleSignClient
}

// New returns a new application instance.
func New(cmd *cobra.Command, version string) (*App, error) {
	config := NewConfiguration()

	a, err := server.New(cmd, server.Options{
		AppName:       "EDIM Mobile Wallet and Issuer API",
		AppVer:        version,
		Configuration: config,
	})
	if err != nil {
		return nil, err
	}

	store, _, err := jsondb.New(a.App, config.Postgres)
	if err != nil {
		return nil, err
	}

	idauth, err := idauth.NewClient(config.IDAuth)
	if err != nil {
		return nil, err
	}

	vci, err := openid4vci.New(a, store, config.Issuer, config.WalletPublicURL)
	if err != nil {
		return nil, err
	}

	att, err := attestation.New(a)
	if err != nil {
		return nil, err
	}

	instance := &App{
		App:         a,
		config:      config,
		store:       store,
		idauth:      idauth,
		vci:         vci,
		attestation: att,
	}

	instance.issuer, err = issuer.NewIssuer(instance, instance.Config().Issuer.TxCodeCacheTTL, instance.Config().WalletPublicURL)
	if err != nil {
		return nil, err
	}

	instance.fprisAPIClient = NewFprisAPIClient(instance.Config().FprisAPIURL)

	instance.rtuAPIClient = NewRTUAPIClient(instance.Config().RTUAPIURL)

	instance.mdlAPIClient = NewMDLAPIClient(instance.Config().MDLAPIURL)

	instance.simpleSignClient, err = NewSimpleSignClient(instance, instance.Config().SimpleSignService, instance.Config().SimpleSignPublicURL, instance.Config().SimpleSignAPIKey, instance.Config().SimpleSignCacheTTL)
	if err != nil {
		return nil, err
	}

	store.AddTask(tasks.NewWalletInstanceCleanupTask(a, store, instance.Config().WalletCheckInterval, instance.Config().WalletOlderThan))

	return instance, nil
}

// Start the application.
func (a *App) Start() error {
	if err := a.Store().Start(a.BackgroundContext()); err != nil {
		return err
	}

	return a.App.Start()
}

// Config returns application configuration.
//
// Panics if configuration is not loaded.
func (a *App) Config() *Configuration {
	if a.config == nil || !a.config.Ready() {
		panic("configuration is not loaded")
	}

	return a.config
}

func (a *App) Store() jsondb.Store {
	return a.store
}

func (a *App) IDAuth() *idauth.Client {
	return a.idauth
}

func (a *App) Issuer() *issuer.Issuer {
	return a.issuer
}

func (a *App) Attestation() *attestation.Service {
	return a.attestation
}

func (a *App) OpenID4VCI() *openid4vci.Service {
	return a.vci
}

func (a *App) FprisAPIClient() *FprisAPIClient {
	return a.fprisAPIClient
}

func (a *App) RTUAPIClient() *RTUAPIClient {
	return a.rtuAPIClient
}

func (a *App) MDLAPIClient() *MDLAPIClient {
	return a.mdlAPIClient
}

func (a *App) SimpleSignClient() *SimpleSignClient {
	return a.simpleSignClient
}
