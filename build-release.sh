#!/usr/bin/env bash
#
# Shell script for building binaries for all relevant platforms
set -euo pipefail

SCRIPT_DIR=$(dirname "$0")
cd "${SCRIPT_DIR}"
DIR_NAME=${PWD##*/} # name of current directory = name of project


# Build
declare -a TARGETS=(darwin linux freebsd openbsd solaris)
VERSION="$(git rev-parse --short HEAD)" # set version to current commit
if git describe --exact-match --tags HEAD 2> /dev/null ; then
    VERSION="$(git describe --exact-match --tags HEAD)"
fi
echo $VERSION

for target in "${TARGETS[@]}" ; do
  output="${DIR_NAME}"
  echo "Building for ${target}, output bin/${output}"
  export CGO_ENABLED=0
  export GOOS=${target}
  export GOARCH=amd64
  go build -ldflags "-s -w" -o "bin/${output}"
  (
    cd ..
    TARBALL="${DIR_NAME}-${VERSION}-${target}-${GOARCH}.tar.gz"
    tar -cf "${TARBALL}" --exclude=.git -z "${DIR_NAME}"
    echo "Created: ${PWD}/${TARBALL}"
  )
  rm -rf "bin/${output}"
done

rmdir bin