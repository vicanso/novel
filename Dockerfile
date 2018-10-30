FROM node:alpine as assets

ADD ./ /novel

RUN apk update \
  && apk add git \
  && git clone --depth=1 https://github.com/vicanso/novel-web /novel/web \
  && cd /novel/admin \
  && yarn \
  && yarn build \
  && cd /novel/web \
  && yarn \
  && yarn build

FROM golang:1.11-alpine as builder

ADD ./ /go/src/github.com/vicanso/novel

COPY --from=assets /novel/admin/dist /go/src/github.com/vicanso/novel/admin/dist
COPY --from=assets /novel/web/dist /go/src/github.com/vicanso/novel/assets

RUN apk update \
  && apk add git \
  && go get -u github.com/golang/dep/cmd/dep \
  && go get -u github.com/gobuffalo/packr/packr \
  && cd /go/src/github.com/vicanso/novel \
  && dep ensure \
  && packr -z \
  && GOOS=linux GOARCH=amd64 go build -tags netgo -o novel

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/vicanso/novel/novel /usr/local/bin/novel
COPY --from=builder /go/src/github.com/vicanso/novel/configs /configs

CMD [ "novel" ]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD novel --check=true || exit 1