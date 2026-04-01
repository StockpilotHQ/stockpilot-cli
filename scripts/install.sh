#!/bin/sh
set -e

REPO="StockpilotHQ/stockpilot-cli"
BINARY="stockpilot"
INSTALL_DIR="/usr/local/bin"

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: $OS"
    echo "On Windows, download the binary from https://github.com/$REPO/releases"
    exit 1
    ;;
esac

# Get latest release version
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
if [ -z "$VERSION" ]; then
  echo "Failed to fetch latest version"
  exit 1
fi

FILENAME="${BINARY}_${OS}_${ARCH}"
URL="https://github.com/$REPO/releases/download/$VERSION/${FILENAME}"

echo "Installing stockpilot $VERSION for ${OS}/${ARCH}..."

TMP=$(mktemp)
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$INSTALL_DIR/$BINARY"
else
  sudo mv "$TMP" "$INSTALL_DIR/$BINARY"
fi

echo "Installed to $INSTALL_DIR/$BINARY"
echo ""
echo "Get started:"
echo "  stockpilot login"
echo "  stockpilot whoami"
echo "  stockpilot status"
