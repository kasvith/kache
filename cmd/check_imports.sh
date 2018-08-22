#!/usr/bin/env bash

set -e

goimports_files=$(goimports -l -local='github.com/kasvith/kache' `find . -name '*.go' | grep -v vendor`)
if [[ -n ${goimports_files} ]]; then
    echo 'goimports needs running on the following files:'
    echo "${goimports_files}"
    exit 1
fi

exit 0