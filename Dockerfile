FROM golang:1.26.6-alpine@sha256:af8d6740070b8906d12eae1c3e3ea0957fb63f492051ea05e354c38ef9fe88df AS builder

COPY --from=ghcr.io/luzifer-docker/pnpm:v11.21.0@sha256:18f46586f69d80207ef2a04f6713feb3b4697c4b882ac80599fe35665ee32a38 . /

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
