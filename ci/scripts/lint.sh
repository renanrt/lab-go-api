#!/bin/sh

go get github.com/alecthomas/gometalinter
gometalinter --install
LINTS="$LINTS deadcode"
LINTS="$LINTS errcheck"
LINTS="$LINTS gofmt"
LINTS="$LINTS gosimple"
LINTS="$LINTS ineffassign"
LINTS="$LINTS unconvert"
LINTS="$LINTS varcheck"
LINTS="$LINTS vet"
gometalinter --vendor --disable-all --deadline=15s $(echo $LINTS|xargs -n 1 echo -n " --enable")
