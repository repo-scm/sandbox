#!/bin/bash

buildTime=$(date +%FT%T%z)
commitID=`git rev-parse --short=7 HEAD`
ldflags="-s -w -X github.com/repo-scm/sandbox/cmd.BuildTime=$buildTime -X github.com/repo-scm/sandbox/cmd.CommitID=$commitID"
target="sandbox"

go env -w GOPROXY=https://goproxy.cn,direct

CGO_ENABLED=0 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target .
