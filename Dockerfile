FROM golang:1.21-alpine@sha256:445f34008a77b0b98bf1821bf7ef5e37bb63cc42d22ee7c21cc17041070d134f as builder

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

