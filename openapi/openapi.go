// SPDX-License-Identifier: EUPL-1.2

//go:generate go run generator/generate.go

// @version 1.0
// @title Auth API
// @description REST API for Authorization and Authentication
// @contactName SIA ZZ Dats
// @contactEmail zzdats@zzdats.lv
// @contactURL https://www.zzdats.lv/
// @server {{SERVER_URL}}
// @security AuthorizationHeader
// @securityScheme AuthorizationHeader http bearer Session identifier. Example: "Authorization: Bearer {session-id}"
package openapi

import (
	_ "embed"
)

//go:embed openapi.json
var OpenAPIDefinition []byte
