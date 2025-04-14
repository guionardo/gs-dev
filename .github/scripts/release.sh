# https://akrabat.com/building-go-binaries-for-different-platforms/

version=$(git describe --tags HEAD)

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    os=$GOOS
    if [ $os = "darwin" ]; then
        os="macOS"
    fi

    output_name="gs-dev-${version}-${os}-${GOARCH}"
    if [ $os = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building release/$output_name..."

    BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

    BUILD_INFO="$BUILD_DATE | ${GOARCH}/${GOOS}"

    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-X 'github.com/guionardo/gs-dev/app/build.BuildInfo=$BUILD_INFO' -X 'github.com/guionardo/gs-dev/app/build.Version=$version'" \
        -o release/$output_name \
        cmd/gs-dev/gs-dev.go

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting.'
        exit 1
    fi
done
