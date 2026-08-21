BUILD_TIME=`date '+%Y-%m-%d %H:%M:%S %z'`
GO_VERSION=`go version`
prefix='duck/lazy.'
LDFLAGS="-X '${prefix}BUILD_TIME=${BUILD_TIME}' -X '${prefix}GO_VERSION=${GO_VERSION}' -X '${prefix}AUTHOR=${AUTHOR}' -s -w"

# Only build services that actually exist
services="admin gamecenter pggateway"

for service in $services; do
    echo "Building $service..."

    if go build -ldflags "${LDFLAGS}" -o bin/$service ./service/$service; then
        echo "✓ Build successful: $service"
    else
        echo "✗ Build failed: $service"
    fi
done

echo ""
echo "Build complete. Binaries available in bin/ directory"
