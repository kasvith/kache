#!/usr/bin/env bash

set -e

fmt_status=$(gofmt -d `find . -name '*.go' | grep -v vendor`)

if ["$fmt_status" == ""]; then
    exit 0
else
    exit 1
fi