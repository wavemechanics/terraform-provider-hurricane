#!/bin/sh
#
# get-binary-semver.sh -- extract semver from binary
#
# The purpose of this script is to produce a filename suffix for the terraform
# provider binary.
# So it uses ident to extract the semver tag and strips any pre-release
# part of it.
#
# So v1.2.3-15 becomes just v1.2.3
#
# Output of ident looks like this:
#
# terraform-provider-hurricane:
#      $semver: v1.2.3 $
#

ident "$1" | awk '
{
    if (NR == 2 && $1 == "$semver:") {
        split($2, a, "-")
        print a[1]
    }
}'
