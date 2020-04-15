from golang:1.13-alpine as builder

add . /go/src/github.com/chentanyi/fileserver
run apk update && apk add git && \
    go get -v github.com/go-bindata/go-bindata/go-bindata && \
    cd /go/src/github.com/chentanyi/fileserver && \
    cd server && go-bindata -pkg server template/ && cd .. && \
    CGO_ENABLED=0 go install

from alpine:latest
workdir /
copy --from=builder /go/bin/fileserver /usr/bin/
entrypoint ["fileserver"]