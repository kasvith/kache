#!/usr/bin/env bash

set -e

# install gimports if not exist
if ! [ -x "$(command -v goimports)" ]; then
    echo 'installing goimports'
    go get -u golang.org/x/tools/cmd/goimports
fi

# check
goimports_files=$(goimports -l -local='github.com/kasvith/kache' `find . -name '*.go' | grep -v vendor`)
if [[ -n ${goimports_files} ]]; then
    echo 'goimports needs running on the following files:'
    echo "${goimports_files}"
    exit 1
fi

exit 0
