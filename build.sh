#!/bin/sh

buildDate=`date +"%F %T %z"`
goVersion=`go version`
goVersion=${goVersion#"go version "}

go build --ldflags "-X main.Version=v0.2.2 \
    -X main.GitCommit=$(git rev-parse HEAD) \
    -X 'main.BuidDate=$buildDate' \
    -X 'main.GoVersion=$goVersion'" \
    -o ./build/net_watcher ./cmd/net_watcher/
