#!/usr/bin/env bash

# Convenience script which just ensures all go code can compile

set -ex
cd "$(dirname "${BASH_SOURCE[0]}")"/../..

export GOBIN=${PWD}/.bin
go list github.com/khulnasoft/khulnasoft/cmd/... |
  grep -o 'github.com/khulnasoft/khulnasoft/cmd/[^/]*$' |
  xargs go install -v

find . \( -name vendor -type d -prune \) -or -name '*_test.go' -exec dirname '{}' \; |
  sort | uniq |
  xargs -n1 -P 8 go test -v -i -c -o /dev/null
