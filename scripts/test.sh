#!/usr/bin/env bash
set -e
echo "" > coverage.txt

export BASE_URL=http://test

for d in $(go list ./... | grep -v vendor); do
    go build
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
