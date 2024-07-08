#!/usr/bin/env bash

function main {
    local release_version=${1:-}

    if [[ -z "${release_version}" ]]; then
        echo "Usage ${0} <release-version>"
        exit 1
    fi

    if ! [[ "${release_version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Release version must be of the form 'v1.2.3'"
        exit 1
    fi

    git commit -am "Release ${release_version}" --allow-empty
    git cliff --tag "${release_version}" -o
    git commit --amend -a
    git tag "${release_version}"
}

main "${@}"
