#!/bin/sh

curdir=`pwd`

cd $(dirname $0)

if [[ -d "server/template" ]]; then
    go get -v github.com/go-bindata/go-bindata/go-bindata
    cd server
    go-bindata -pkg server template/
    cd ..
fi

cd $curdir