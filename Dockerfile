FROM golang:1.27.3-alpine@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.28.11@sha256:15f1cb193c684b105a6410d8dde0ac13fba77f100315311bfc7775b720b194d5 . /

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
