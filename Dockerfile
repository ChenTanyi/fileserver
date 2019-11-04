from golang:1.13-alpine as builder

add . /go/src/github.com/chentanyi/fileserver
run cd /go/src/github.com/chentanyi/fileserver && \
    CGO_ENABLED=0 go install

from alpine:latest
workdir /app
copy --from=builder /go/bin/fileserver .
entrypoint ["./fileserver"]