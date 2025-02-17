#!/usr/bin/env bash
echo "Building and installing the application..."

BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git log -n 1 "--pretty=format:%h %ae %s")"
VERSION="$(git describe --tags --abbrev=0 | tr -d '\n')"
HOSTNAME=$(hostname)
HAS_CHANGES=$(git diff --name-only HEAD | wc -l | xargs)
if [[ $HAS_CHANGES -gt 0 ]]; then
  VERSION="$VERSION-dev-$HAS_CHANGES"
fi
echo VERSION: $VERSION

_goenv="$(go env -json)"
_goarch="$(echo $_goenv | jq .GOARCH | tr -d \")"
_goos="$(echo $_goenv | jq .GOOS | tr -d \")"
BUILD_INFO="$BUILD_DATE | ${USER}@${HOSTNAME} | ${_goarch}/${_goos} | ${GIT_COMMIT}"
echo BUILD_INFO: $BUILD_INFO

if [ "$1" == "install" ]; then
  mkdir -p bin
  cmd="build -o bin/gs-dev"
  oper="building to bin/gs-dev"
else
  cmd="install"
  oper="installing on $GOPATH/bin/gs-dev"
fi

go $cmd -ldflags="-X 'github.com/guionardo/gs-dev/app/build.BuildInfo=$BUILD_INFO' -X 'github.com/guionardo/gs-dev/app/build.Version=$VERSION'" cmd/gs-dev.go
if [ $? == 0 ]; then
  echo "Operation: $oper OK"
else
  echo "Operation: $oper FAILED"
  exit 1
fi
