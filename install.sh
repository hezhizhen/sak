#!/bin/sh
set -e

REPO="hezhizhen/sak"
BINARY="sak"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"([^"]+)".*/\1/'
}

detect_os() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)  echo "linux" ;;
        darwin) echo "darwin" ;;
        *)      echo "Unsupported OS: $os" >&2; exit 1 ;;
    esac
}

detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac
}

VERSION="${1:-$(get_latest_version)}"
if [ -z "$VERSION" ]; then
    echo "Error: could not determine version" >&2
    exit 1
fi

OS=$(detect_os)
ARCH=$(detect_arch)
# Strip leading 'v' for filename
FILE_VERSION="${VERSION#v}"

FILENAME="${BINARY}_${FILE_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${BINARY} ${VERSION} (${OS}/${ARCH})..."
curl -sL "$URL" -o "${TMPDIR}/${FILENAME}"

tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR"

echo "Installing to ${INSTALL_DIR}/${BINARY}..."
install -m 755 "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

echo "Done: $(${BINARY} version --short)"
