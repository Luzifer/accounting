FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.28.5@sha256:b87bd3ca3ba0d31c0ce5c87eb00b1b2cb7f03649502e35437f6c548ef788961b . /

COPY . /go/src/accounting
WORKDIR /go/src/accounting

RUN set -ex \
 && apk add --no-cache \
      git \
      make \
      nodejs-current \
 && make frontend build


FROM scratch

LABEL org.opencontainers.image.authors='Knut Ahlers <knut@ahlers.me>' \
      org.opencontainers.image.url='https://github.com/Luzifer/accounting/pkgs/container/accounting' \
      org.opencontainers.image.documentation='https://github.com/Luzifer/accounting' \
      org.opencontainers.image.source='https://github.com/Luzifer/accounting' \
      org.opencontainers.image.licenses='Apache-2.0'

COPY --from=builder /go/src/accounting/accounting /usr/bin/accounting

EXPOSE 3000
USER 1000:1000

ENTRYPOINT ["/usr/bin/accounting"]
