FROM golang:1.26-alpine@sha256:c2a1f7b2095d046ae14b286b18413a05bb82c9bca9b25fe7ff5efef0f0826166 AS builder

COPY . /go/src/accounting
WORKDIR /go/src/accounting

RUN set -ex \
 && apk add --no-cache \
      git \
      make \
      nodejs \
      npm \
 && make frontend build


FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/src/accounting/accounting /usr/local/bin/accounting

EXPOSE 3000
USER 1000

ENTRYPOINT ["/usr/local/bin/accounting"]
CMD ["--"]

# vim: set ft=Dockerfile:
