#!/bin/sh

mkdir "$GOPATH"/src/github.com/roman-mazur
cd "$GOPATH"/src/github.com/roman-mazur || echo "Couldn't open dir $GOPATH/src/github.com/roman-mazur"

go get -u github.com/roman-mazur/bood/

cd ./bood/ || echo "Couldn't open dir bood"
go run ./cmd/bood/main.go

export PATH=$(pwd)/out/bin:$PATH
