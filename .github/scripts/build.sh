#!/usr/bin/env bash
echo "Building and installing the application..."
mkdir -p bin
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git log -n 1 "--pretty=format:%h %ae %s")"
VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')"
HOSTNAME=$(hostname)
HAS_CHANGES=$(git diff --name-only HEAD | wc -l | xargs)
if [[ $HAS_CHANGES -gt 0 ]]; then
  VERSION="$VERSION-dev-$HAS_CHANGES"
fi
echo VERSION: $VERSION
BUILD_INFO="$BUILD_DATE | ${USER}@${HOSTNAME} | ${GIT_COMMIT}"
echo BUILD_INFO: $BUILD_INFO
go build -o bin/gs-dev -ldflags="-X 'github.com/guionardo/gs-dev/app/build.BuildInfo=$BUILD_INFO' -X 'github.com/guionardo/gs-dev/app/build.Version=$VERSION'" cmd/cli.go
echo "Application built successfully -> bin/gs-dev"
