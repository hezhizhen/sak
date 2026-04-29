#!/bin/sh
set -e

REPO="hezhizhen/sak"
BINARY="sak"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"

get_latest_version() {
    if command -v gh >/dev/null 2>&1; then
        gh release view --repo "${REPO}" --json tagName -q .tagName 2>/dev/null && return
    fi
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"([^"]+)".*/\1/'
}

detect_os() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        linux)  echo "linux" ;;
        darwin) echo "darwin" ;;
        *)      echo "Error: unsupported OS \"$os\". Only linux and darwin are supported." >&2; exit 1 ;;
    esac
}

detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             echo "Error: unsupported architecture \"$arch\". Only amd64 and arm64 are supported." >&2; exit 1 ;;
    esac
}

VERSION="${1:-$(get_latest_version)}"
if [ -z "$VERSION" ]; then
    echo "Error: could not determine latest version." >&2
    echo "Please specify a version manually: sh install.sh v0.1.0" >&2
    echo "Or install gh CLI for authenticated GitHub API access: https://cli.github.com" >&2
    exit 1
fi

OS=$(detect_os)
ARCH=$(detect_arch)
FILE_VERSION="${VERSION#v}"

FILENAME="${BINARY}_${FILE_VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${BINARY} ${VERSION} (${OS}/${ARCH})..."
HTTP_CODE=$(curl -sL -o "${TMPDIR}/${FILENAME}" -w "%{http_code}" "$URL")
if [ "$HTTP_CODE" != "200" ] || [ ! -s "${TMPDIR}/${FILENAME}" ]; then
    echo "Error: failed to download ${FILENAME} (HTTP ${HTTP_CODE})." >&2
    echo "Check that version ${VERSION} exists: https://github.com/${REPO}/releases" >&2
    exit 1
fi

if ! tar -xzf "${TMPDIR}/${FILENAME}" -C "$TMPDIR" 2>/dev/null; then
    echo "Error: failed to extract ${FILENAME}. The downloaded file may be corrupted." >&2
    exit 1
fi

if [ ! -f "${TMPDIR}/${BINARY}" ]; then
    echo "Error: binary \"${BINARY}\" not found in archive." >&2
    exit 1
fi

echo "Installing to ${INSTALL_DIR}/${BINARY}..."
mkdir -p "${INSTALL_DIR}"
if ! install -m 755 "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"; then
    echo "Error: failed to install to ${INSTALL_DIR}. Check directory permissions." >&2
    exit 1
fi

echo "Done: ${BINARY} ${VERSION} installed to ${INSTALL_DIR}/${BINARY}"
