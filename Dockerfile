FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk --no-cache add git ca-certificates tzdata

COPY go.mod go.sum ./

COPY . ./

RUN go generate ./...

RUN go build -ldflags="-w -s" -tags 'netgo osusergo' -o publish/server ./cmd/server/... 

RUN mkdir -p publish/etc/ssl/certs/ && \
    mkdir -p publish/usr/share/zoneinfo/ && \
    mkdir -p publish/certs/ && \
    cp /etc/ssl/certs/ca-certificates.crt publish/etc/ssl/certs/ && \
    cp -R /usr/share/zoneinfo publish/usr/share/

FROM scratch
WORKDIR /
COPY --from=build app/publish/ ./
EXPOSE 8080/tcp
ENV TZ=Europe/Riga

ENTRYPOINT ["/server", "main"]
HEALTHCHECK --start-period=30s --start-interval=5s --interval=1m --timeout=10s --retries=5 CMD ["/server", "health"]