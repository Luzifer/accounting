FROM golang:1.27.3-alpine@sha256:f2c872bcee949b7be9c490efd02b48a8a2fbf947cdf533f1c58525ec9d8170fd AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.28.6@sha256:d35ec722ba2d0bc7d8ff9d9920d2b9f851e566d13822e33da1272d085a5b03a8 . /

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
