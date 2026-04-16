#!/usr/bin/env sh
# install.sh — one-line installer for balatro-cli.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/daghlny/balatro-cli/main/install.sh | sh
#
# Environment overrides:
#   INSTALL_DIR   install location (default: $HOME/.local/bin)
#   VERSION       release tag to install (default: latest release)
set -eu

REPO="daghlny/balatro-cli"
BINARY="balatro-cli"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${VERSION:-}"

# --- detect OS / arch ----------------------------------------------------
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)    arch=amd64 ;;
  aarch64|arm64)   arch=arm64 ;;
  *) echo "Unsupported arch: $arch" >&2; exit 1 ;;
esac
case "$os" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $os (use the zip from the Releases page on Windows)" >&2; exit 1 ;;
esac

# --- resolve version -----------------------------------------------------
if [ -z "$VERSION" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n1)
fi
if [ -z "$VERSION" ]; then
  echo "Could not resolve a release tag. Has a release been published yet?" >&2
  echo "See https://github.com/$REPO/releases" >&2
  exit 1
fi
semver=${VERSION#v}

archive="${BINARY}_${semver}_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$VERSION/$archive"

# --- download + install --------------------------------------------------
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $archive ($VERSION)..."
if ! curl -fsSL "$url" -o "$tmp/$archive"; then
  echo "Download failed. URL: $url" >&2
  exit 1
fi

tar -xzf "$tmp/$archive" -C "$tmp"
if [ ! -f "$tmp/$BINARY" ]; then
  echo "Archive did not contain $BINARY. Contents:" >&2
  ls -la "$tmp" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"

echo
echo "Installed $BINARY $VERSION -> $INSTALL_DIR/$BINARY"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Note: $INSTALL_DIR is not on your PATH. Add this to your shell rc:"
     echo "      export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac
