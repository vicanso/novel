FROM golang:1.11-alpine as builder

ADD ./ /go/src/github.com/vicanso/novel

RUN apk update \
  && apk add git \
  && go get -u github.com/golang/dep/cmd/dep \
  && cd /go/src/github.com/vicanso/novel \
  && dep ensure \
  && GOOS=linux GOARCH=amd64 go build -tags netgo -o novel

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/vicanso/novel/novel /usr/local/bin/novel

CMD [ "novel" ]