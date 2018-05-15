FROM golang:alpine

ARG SOURCE_COMMIT

ADD . /go/src/github.com/vdemay/sharedcount-exporter
WORKDIR /go/src/github.com/vdemay/sharedcount-exporter

RUN DATE=$(date -u '+%Y-%m-%d-%H%M UTC'); \
    go install -ldflags="-X 'main.Version=${SOURCE_COMMIT}' -X 'main.BuildTime=${DATE}'" ./...

ENTRYPOINT  [ "/go/bin/sharedcount-exporter" ]
EXPOSE      9383
