FROM golang:alpine as builder

COPY . /go/src/accounting
WORKDIR /go/src/accounting

RUN set -ex \
 && apk add --no-cache \
      git \
      make \
      nodejs \
      npm \
 && make build


FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/src/accounting/accounting /usr/local/bin/accounting

EXPOSE 3000
USER 1000

ENTRYPOINT ["/usr/local/bin/accounting"]
CMD ["--"]

# vim: set ft=Dockerfile:
