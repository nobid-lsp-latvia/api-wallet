// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"time"

	"git.zzdats.lv/edim/api-wallet/issuer"
	jsondb "github.com/nobid-lsp-latvia/lx-go-jsondb"

	"azugo.io/azugo/config"
	corecfg "azugo.io/core/config"
	"azugo.io/core/validation"
	"github.com/nobid-lsp-latvia/go-idauth"
	"github.com/spf13/viper"
)

// Configuration represents the configuration for the application.
type Configuration struct {
	*config.Configuration `mapstructure:",squash"`

	Postgres *jsondb.Configuration `mapstructure:"postgres"`
	IDAuth   *idauth.Configuration `mapstruct:"idauth"`
	Issuer   *issuer.Configuration `mapstruct:"issuer"`

	QRAPIDeepLink       string        `mapstructure:"qr_api_deep_link" validate:"required"`
	FprisAPIURL         string        `mapstructure:"fpris_api_url" validate:"required,url"`
	RTUAPIURL           string        `mapstructure:"rtu_api_url" validate:"required,url"`
	MDLAPIURL           string        `mapstructure:"mdl_api_url" validate:"required,url"`
	SimpleSignService   string        `mapstructure:"simple_sign_service" validate:"required,url"`
	SimpleSignPublicURL string        `mapstructure:"simple_sign_public_url" validate:"required,url"`
	SimpleSignAPIKey    string        `mapstructure:"simple_sign_api_key" validate:"required"`
	SimpleSignCacheTTL  time.Duration `mapstructure:"simple_sign_cache_ttl" validate:"required,gt=0"`
	WalletCheckInterval time.Duration `mapstructure:"wallet_check_interval" validate:"required"`
	WalletOlderThan     time.Duration `mapstructure:"wallet_older_than" validate:"required"`
	WalletPublicURL     string        `mapstructure:"wallet_api_public_url" validate:"required,url"`
}

// NewConfiguration returns a new configuration.
func NewConfiguration() *Configuration {
	return &Configuration{
		Configuration: config.New(),
	}
}

// Core returns the core configuration.
func (c *Configuration) ServerCore() *config.Configuration {
	return c.Configuration
}

// Bind configuration to viper.
func (c *Configuration) Bind(_ string, v *viper.Viper) {
	c.Configuration.Bind("", v)

	c.Postgres = config.Bind(c.Postgres, "postgres", v)
	c.Issuer = config.Bind(c.Issuer, "issuer", v)
	c.IDAuth = config.Bind(c.IDAuth, "idauth", v)

	v.SetDefault("wallet_check_interval", 30*time.Minute)
	v.SetDefault("wallet_older_than", 1*time.Hour)
	v.SetDefault("simple_sign_cache_ttl", 10*time.Minute)

	_ = v.BindEnv("qr_api_deep_link", "QR_API_DEEP_LINK")
	_ = v.BindEnv("fpris_api_url", "FPRIS_API_URL")
	_ = v.BindEnv("rtu_api_url", "RTU_API_URL")
	_ = v.BindEnv("mdl_api_url", "MDL_API_URL")
	_ = v.BindEnv("wallet_api_public_url", "WALLET_API_PUBLIC_URL")
	_ = v.BindEnv("simple_sign_service", "SIMPLE_SIGN_SERVICE")
	_ = v.BindEnv("simple_sign_public_url", "SIMPLE_SIGN_PUBLIC_URL")
	_ = v.BindEnv("simple_sign_cache_ttl", "SIMPLE_SIGN_CACHE_TTL")

	key, _ := corecfg.LoadRemoteSecret("SIMPLE_SIGN_API_KEY")
	v.SetDefault("simple_sign_api_key", key)
	_ = v.BindEnv("simple_sign_api_key", "SIMPLE_SIGN_API_KEY")
	_ = v.BindEnv("wallet_check_interval", "WALLET_CHECK_INTERVAL")
	_ = v.BindEnv("wallet_older_than", "WALLET_OLDER_THAN")
}

// Validate application configuration.
func (c *Configuration) Validate(validate *validation.Validate) error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	if err := c.Postgres.Validate(validate); err != nil {
		return err
	}

	if err := c.Issuer.Validate(validate); err != nil {
		return err
	}

	if err := c.IDAuth.Validate(validate); err != nil {
		return err
	}

	return nil
}
