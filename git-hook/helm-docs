#!/usr/bin/env bash

set -e

if ! command -v helm-docs > /dev/null 2>&1; then
    echo "Please install helm-docs to run the pre-commit hook! https://github.com/norwoodj/helm-docs#installation"
    exit 1
fi

helm-docs "${@}"
