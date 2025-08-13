// SPDX-License-Identifier: EUPL-1.2

package openid4vci

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"git.zzdats.lv/edim/api-wallet/issuer"
	jsondb "github.com/nobid-lsp-latvia/lx-go-jsondb"

	"aidanwoods.dev/go-paseto"
	"azugo.io/azugo"
	"azugo.io/core/cache"
)

type Service struct {
	app    *azugo.App
	config *issuer.Configuration
	store  jsondb.Store

	nonceCache cache.Instance[bool]
	nonceLock  sync.Mutex
	nonceKey   paseto.V4SymmetricKey

	walletPublicURL   string
	walletInstanceURL string
}

func New(app *azugo.App, store jsondb.Store, config *issuer.Configuration, publicBaseURL string) (*Service, error) {
	b, err := base64.StdEncoding.DecodeString(config.NonceSharedSecret)
	if err != nil {
		return nil, err
	}

	key, err := paseto.V4SymmetricKeyFromBytes(b)
	if err != nil {
		return nil, err
	}

	cache, err := cache.Create[bool](app.Cache(), "nonce-reuse", cache.DefaultTTL(config.NonceTTL))
	if err != nil {
		return nil, err
	}

	walletInstanceURL, err := url.JoinPath(publicBaseURL, "instance")
	if err != nil {
		return nil, fmt.Errorf("failed to generate instance ID: %w", err)
	}

	return &Service{
		app:    app,
		config: config,
		store:  store,

		nonceCache: cache,
		nonceKey:   key,

		walletPublicURL:   strings.TrimSuffix(publicBaseURL, "/"),
		walletInstanceURL: walletInstanceURL,
	}, nil
}
