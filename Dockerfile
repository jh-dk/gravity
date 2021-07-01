FROM alpine:3.13 AS builder

RUN set -ex && \
	apk add --no-cache make bash libc6-compat aws-cli fakeroot

WORKDIR /go/src/github.com/gravitational/gravity
