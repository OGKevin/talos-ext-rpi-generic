FROM golang:1.20-alpine as builder

WORKDIR /app

RUN apk add ca-certificates curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d 

# ENV USER=generic-pi
# ENV UID=10001
#
# RUN adduser \
#   --disabled-password \
#   --gecos "" \
#   --home "/nonexistent" \
#   --shell "/sbin/nologin" \
#   --no-create-home \
#   --uid "${UID}" \
#   "${USER}"
#
COPY go.* /app/

RUN --mount=type=cache,target=$GOPATH/pkg go mod download

COPY . /app/

RUN  ./bin/task build-binary

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/_out/generic-pi /
# COPY --from=builder /etc/passwd /etc/passwd

# USER generic-pi

ENTRYPOINT ["/generic-pi"]
