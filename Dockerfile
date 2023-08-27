FROM golang:1.20-alpine@sha256:03278bc16e1a5b4fb6cdd3462108c060aa1e9c2353ce4d15d744b3c40168677d as builder

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

