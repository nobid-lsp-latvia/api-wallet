# Api-Wallet

EDIM Mobile Wallet and Issuer API

## Built with Azugo Go Web Framework

This project is built using the [Azugo Go Web Framework](https://azugo.io), a powerful and flexible framework for building modern web applications in Go. Check out the [Azugo GitHub page](https://github.com/azugo) for more information and documentation.

<!-- TOC -->

- [Api-Wallet](#repo_name_title)
  - [Development](#development)
    - [Prepare dependecies](#prepare-dependecies)
    - [Local development](#local-development)
    - [Before commit](#before-commit)
  - [Environment variables](#environment-variables)
    - [Local example](#local-example)

<!-- /TOC -->

## Development

### Prepare dependecies

```sh
go mod download
go generate ./...
```

### Local development

To build in VS Code use `Ctrl`+`Shift`+`B`.

To debug project in VS Code use `F5`.

### Before commit

> CI requires linted, formatted code

You should run:

```sh
gofmt -s -w ./..
```

or

```sh
gofumpt -w ./..
```

and fix any errors reported by

```sh
golangci-lint run
```

## Environment variables

In order to run the service you need configure environment variables. List of environment variables:

| Variable | Description | Default value | Required |
| --- | --- | --- | --- |
| `SERVER_URLS` | An server URL or multiple URLS separated by semicolon to listen on. | 0.0.0.0:8080 | Yes |
| `ENVIRONMENT` | Environment name. Possible values: `Development`, `Staging`, `Production` | `Development` | Yes |
| `BASE_PATH` | Base path for all routes | `/` (or take value from `SERVER_URLS` path if exists) | No |
| `ACCESS_LOG_ENABLED` | Enable access log | `true` | Yes |
| `REVERSE_PROXY_LIMIT` | Limit for reverse proxy. | `1` | No |
| `REVERSE_PROXY_TRUSTED_IPS` | List of trusted IP addresses for reverse proxy. Separated by `;` | `"127.0.0.1"` | No |
| `REVERSE_PROXY_TRUSTED_HEADERS` | List of trusted headers for reverse proxy. Separated by `;` | `X-Real-IP; X-Forwarded-For` | No |
| `LOG_LEVEL` | Minimal log level. Allowed values are `debug`, `info`, `warn`, `error`, `fatal`, `panic` | `info` | Yes |
| `CACHE_TYPE` | Cache type to use in service. Allowed values are `memory`, `redis`, `redis-cluster`. | `memory` | No |
| `CACHE_TTL` | Duration on how long to keep items in cache. Defaults to 0 meaning to never expire. | `0` | No |
| `CACHE_KEY_PREFIX` | Prefix all cache keys with specified value. | `""` | No |
| `CACHE_CONNECTION` | If other than memory cache is used specifies connection string on how to connect to cache storage. | `""` | No |
| `CACHE_PASSWORD` / `CACHE_PASSWORD_FILE` | Password to use in connection string. | `""` | No |
| `POSTGRES_HOST` | PostgreSQL HOST FQDN | `"db.example.lv"` | Yes |
| `POSTGRES_PORT` | PostgreSQL port | `"5432"` | Yes |
| `POSTGRES_USER` | PostgreSQL  | `"edim_public"` | Yes |
| `POSTGRES_DB` | PostgreSQL  | `"edim"` | Yes |
| `POSTGRES_PASSWORD` | PostgreSQL  | `/secret/edim-public-db-pw` | Yes |
| `IDAUTH_URL` | "" | URL for IDAuth service (empty/not configured) | yes |
| `IDAUTH_CLIENT_ID` | "" | `api-wallet` id registrated in idAuth service | yes |
| `IDAUTH_CLIENT_SECRET_FILE` | "/secret/edim-idauth-client-secret-api-mdl-data" | Path to the file containing the client secret for authentication | Yes |
| `ISSUER_NONCE_SHARED_SECRET` / `ISSUER_NONCE_SHARED_SECRET_FILE` | Base64 encoded 256-bit secret that is used to encrypt and decrypt nonce | `""` | Yes |
| `ISSUER_API_URL` | Internal URL for the `demo-issuer` service | `"http://demo-issuer.edim-test.svc.cluster.local:5000"` | Yes |
| `ISSUER_CERTIFICATE_FILE` | Path to issuer signing certificate PEM file | `/secret/edim-issuer-certificate` | Yes |
| `ISSUER_CERTIFICATE_PASSWORD` / `ISSUER_CERTIFICATE_PASSWORD_FILE` | Issuer signing certificate PEM password | `""` | No |
| `ISSUER_API_URL` | Internal URL for the `demo-issuer` service | `"http://demo-issuer.edim-test.svc.cluster.local:5000"` | Yes |
| `FPRIS_API_URL` | Internal URL for the `api-fpris` service | `"http://api-fpris.edim-test.svc.cluster.local:8080/fpris"` | Yes |
| `RTU_API_URL` | Internal URL for the `api-rtu` service|`"http://api-rtu.edim-test.svc.cluster.local:8080/rtu"` | Yes |
| `MDL_API_URL` | Internal URL for the `api-mdl` service|`"http://api-mdl.edim-test.svc.cluster.local:8080/mdl"` | Yes |
| `WALLET_API_PUBLIC_URL` | Public URL for the `api-wallet` service, that will be put in deeplink|`"https://edim-api-dev.local/wallet"` | Yes |
| `AUDIT_ENDPOINT` | Internal URL for the `api-audit` service|`"http://api-audit.edim-test.svc.cluster.local:8080/audit/1.0"` | Yes |
| `QR_API_DEEP_LINK` | Schema to add to response|`"openid-credential-offer"` | Yes |
| `SIMPLE_SIGN_SERVICE` | URL for simple sign service|`"https://signing.example.lv/simple-sign/"` | Yes |
| `SIMPLE_SIGN_PUBLIC_URL` | Public URL for simple sign service|`"https://signing.example.lv/simple-sign/"` | Yes |
| `SIMPLE_SIGN_API_KEY` | API key simple sign service|`"examplekey"` | No |
| `SIMPLE_SIGN_CACHE_TTL` | Duration on how long to keep items in cache for simple sign service. Defaults to 10min. | `10m` | No |
| `WALLET_CHECK_INTERVAL` | How frequently check expired wallet instances (e.g., `1m` for 1 min, `1h` for 1 hour, `1d` for 1 day) |`"30m"` | No |
| `WALLET_OLDER_THAN` | Duration specifying how long wallet instances remain valid (e.g., `1m` for 1 min, `1h` for 1 hour, `1d` for 1 day)  |`"1h"` | No |

### Local example

Generate nonce encryptoion secret using command line:

```sh
openssl rand -base64 32
```

In local development you must create `.env` file in the root of the project. Example:

```sh
IDAUTH_URL=https://edim-dev.local/idauth/
IDAUTH_CLIENT_ID=edim-demo
IDAUTH_CLIENT_SECRET=edim-demo

POSTGRES_HOST=edim-db-dev.local
POSTGRES_PORT=5432
POSTGRES_USER=wallet_public
POSTGRES_DB=edim
POSTGRES_PASSWORD=xxx

ISSUER_API_URL=https://edim-demo-issuer-dev.local/
ISSUER_NONCE_SHARED_SECRET=xxx

QR_API_DEEP_LINK=test
FPRIS_API_URL=https://edim-api-dev.local/fpris/
RTU_API_URL=https://edim-api-dev.local/rtu/
MDL_API_URL=https://edim-api-dev.local/mdl/
WALLET_API_PUBLIC_URL=https://edim-api-dev.local/wallet

WALLET_CHECK_INTERVAL=30m
WALLET_OLDER_THAN=1h
```
