// SPDX-License-Identifier: EUPL-1.2

package issuer

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"time"

	"azugo.io/core/cert"
	"azugo.io/core/config"
	"azugo.io/core/validation"
	"github.com/spf13/viper"
)

// Configuration of the Issuer.
type Configuration struct {
	NonceTTL                 time.Duration `mapstructure:"nonce_ttl" validate:"required,gt=0"`
	NonceSharedSecret        string        `mapstructure:"nonce_shared_secret" validate:"required,base64"`
	IssuerCertificate        string        `mapstructure:"issuer_certificate" validate:"required"`
	IssuerCertificatePasword string        `mapstructure:"issuer_certificate_password"`
	APIURL                   string        `mapstructure:"issuer_url" validate:"required"`
	TxCodeCacheTTL           time.Duration `mapstructure:"issuer_tx_cache_ttl" validate:"required,gt=0"`

	signingCertificate *tls.Certificate
}

func (c *Configuration) Bind(prefix string, v *viper.Viper) {
	nonceSharedSecret, _ := config.LoadRemoteSecret("ISSUER_NONCE_SHARED_SECRET")
	v.SetDefault(prefix+".nonce_shared_secret", nonceSharedSecret)
	v.SetDefault(prefix+".nonce_ttl", 10*time.Minute)

	certificate, _ := config.LoadRemoteSecret("ISSUER_CERTIFICATE")
	v.SetDefault(prefix+".issuer_certificate", certificate)

	password, _ := config.LoadRemoteSecret("ISSUER_CERTIFICATE_PASSWORD")
	v.SetDefault(prefix+".issuer_certificate_password", password)

	v.SetDefault(prefix+".issuer_tx_cache_ttl", 10*time.Minute)

	_ = v.BindEnv(prefix+".nonce_shared_secret", "ISSUER_NONCE_SHARED_SECRET")
	_ = v.BindEnv(prefix+".nonce_ttl", "ISSUER_NONCE_TTL")
	_ = v.BindEnv(prefix+".issuer_certificate", "ISSUER_CERTIFICATE")
	_ = v.BindEnv(prefix+".issuer_certificate_password", "ISSUER_CERTIFICATE_PASSWORD")
	_ = v.BindEnv(prefix+".issuer_url", "ISSUER_API_URL")
	_ = v.BindEnv(prefix+".issuer_tx_cache_ttl", "ISSUER_TX_CACHE_TTL")
}

func (c *Configuration) SigningCertificate() (*tls.Certificate, error) {
	// Skip if the certificate is already loaded
	if c.signingCertificate != nil {
		return c.signingCertificate, nil
	}

	certbuf, keybuf, err := cert.LoadPEMFromReader(bytes.NewReader([]byte(c.IssuerCertificate)), cert.Password(c.IssuerCertificatePasword))
	if err != nil {
		return nil, err
	}

	signCert, err := cert.LoadTLSCertificate(certbuf, keybuf)
	if err != nil {
		return nil, err
	}

	c.signingCertificate = signCert

	return c.signingCertificate, nil
}

// Validate Issuer configuration section.
func (c *Configuration) Validate(valid *validation.Validate) error {
	if err := valid.Struct(c); err != nil {
		return err
	}

	b, err := base64.StdEncoding.DecodeString(c.NonceSharedSecret)
	if err != nil {
		return errors.New("nonce_shared_secret must be a valid base64 encoded string")
	}

	if len(b) != 32 {
		return errors.New("nonce_shared_secret must be exactly 32 bytes long")
	}

	_, err = c.SigningCertificate()

	return err
}
