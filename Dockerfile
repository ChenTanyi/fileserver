from golang:1.16-alpine as builder

add . /go/src/github.com/chentanyi/fileserver
run apk update && apk add git && \
    cd /go/src/github.com/chentanyi/fileserver && \
    CGO_ENABLED=0 go get ./... && \
    CGO_ENABLED=0 go install

from alpine:latest
workdir /
copy --from=builder /go/bin/fileserver /usr/bin/
entrypoint ["fileserver"]
