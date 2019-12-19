#!/bin/bash
#
# get-version.sh -- get version based on git tags and state of local repo

if ! test -d .git
then
    echo "$0: not in a git working directory, not changing version" 1>&2
    exit 1
fi

version=$(git describe --tags --always --dirty 2>/dev/null)
if test -z "$version"
then
    echo "$0: problem getting current git tags and version, not changing version" 2>&2
    exit 1
fi

echo $version
