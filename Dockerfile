FROM golang:1.21-alpine@sha256:96634e55b363cb93d39f78fb18aa64abc7f96d372c176660d7b8b6118939d97b as builder

WORKDIR /app

RUN apk add ca-certificates curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d 

COPY go.* /app/

RUN --mount=type=cache,target=$GOPATH/pkg go mod download

COPY . /app/

RUN  ./bin/task build-binary

FROM scratch

LABEL org.opencontainers.image.source https://github.com/OGKevin/talos-ext-rpi-generic

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/_out/generic-pi /

ENTRYPOINT ["/generic-pi"]

