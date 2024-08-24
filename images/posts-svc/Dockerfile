FROM golang@sha256:ce1e987ea7759217351b74977a384cea8f44631f1c1add04d1703f13dd3ee850 AS builder

ENV CGO_ENABLED=0

WORKDIR /go/app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY . .

WORKDIR /go/app/posts
RUN --mount=type=cache,target=/gomod-cache go mod download
RUN --mount=type=cache,target=/go-cache go build -o app

FROM bash

COPY --from=builder /go/app/posts/app /go/app/hs256secret.txt .

ENTRYPOINT ["./app"]
