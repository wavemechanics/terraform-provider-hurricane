#!/bin/bash

usage() {
    cat <<EOF
This script embeds a version in the named file.
It does this by changing the contents of any CVS $semver tag in named file.

usage: $0 <target-file> <new-version>

where:
    target-file is the file with the tag to be changed
    new-version is the version to embed
EOF
    exit $1
}

main() {
    local path=$1
    local version=$2

    if test -z "$path"
    then
        usage 1
    fi
    if test -z "$version"
    then
        echo "$0: empty version, not touching $path"
        exit 0
    fi

    sed -e 's|\(\$semver: \)\(.*\)\( \$\)|\1'"$version"'\3|' < "$path" > "$path".new
    if test $? != 0
    then
        exit
    fi
    mv "$path".new "$path"
}

main "$@"
