// SPDX-License-Identifier: EUPL-1.2

package attestation

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"azugo.io/azugo"
	"github.com/fxamacker/cbor/v2"
)

type Service struct {
	app *azugo.App
}

func New(app *azugo.App) (*Service, error) {
	return &Service{
		app: app,
	}, nil
}

type Result struct {
	HardwareKeyTag string `json:"hardwareKeyTag"`
	CertsIssued    int    `json:"certsIssued,omitempty"`
	DeviceType     string `json:"deviceType"`
	PublicKey      string `json:"publicKey"`
}

func (a *Service) Verify(att string, challenge []byte, tag string) (*Result, error) {
	buf, err := base64.RawURLEncoding.DecodeString(att)
	if err != nil {
		return nil, err
	}

	af := struct {
		Format string `cbor:"fmt"`
	}{}

	decoder := cbor.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(&af); err != nil {
		return nil, err
	}

	switch af.Format {
	case "android-key", "android":
		return a.verifyAndroid(att, challenge, tag, time.Now())
	case "apple-appattest", "apple":
		return a.verifyIOS(att, challenge, tag, time.Now())
	default:
		return nil, fmt.Errorf("unknown attestation format: %s", af.Format)
	}
}
