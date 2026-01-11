#!/bin/bash
set -e

# Subman Release Script
# Usage: ./release.sh v1.0.3

if [ -z "$1" ]; then
  echo "Usage: ./release.sh <version>"
  echo "Example: ./release.sh v1.0.3"
  exit 1
fi

VERSION=$1

echo "🚀 Building Subman $VERSION for all platforms..."
echo ""

# Check Docker is running
if ! docker ps &> /dev/null; then
  echo "❌ Docker is not running. Please start Docker and try again."
  exit 1
fi

# Build all platforms
echo "📦 Building macOS (Intel + Apple Silicon)..."
fyne-cross darwin -arch=amd64,arm64 -app-id=com.subman.app

echo ""
echo "📦 Building Linux (amd64 + arm64)..."
fyne-cross linux -arch=amd64,arm64 -app-id=com.subman.app

echo ""
echo "📦 Building Windows (amd64 + arm64)..."
fyne-cross windows -arch=amd64,arm64 -app-id=com.subman.app

# Package macOS .app bundles
echo ""
echo "📦 Packaging macOS .app bundles..."
cd fyne-cross/dist/darwin-amd64
zip -r -q ../subman-macos-amd64.zip subman.app
cd ../../..

cd fyne-cross/dist/darwin-arm64
zip -r -q ../subman-macos-arm64.zip subman.app
cd ../../..

# Copy and rename binaries with descriptive names
echo ""
echo "📦 Preparing release assets with unique names..."
cp fyne-cross/dist/linux-amd64/subman.tar.xz fyne-cross/dist/subman-linux-amd64.tar.xz
cp fyne-cross/dist/linux-arm64/subman.tar.xz fyne-cross/dist/subman-linux-arm64.tar.xz
cp fyne-cross/dist/windows-amd64/subman.exe.zip fyne-cross/dist/subman-windows-amd64.zip
cp fyne-cross/dist/windows-arm64/subman.exe.zip fyne-cross/dist/subman-windows-arm64.zip

echo ""
echo "✅ All binaries built successfully!"
echo ""
echo "📋 Built artifacts:"
ls -lh fyne-cross/dist/subman-macos-amd64.zip
ls -lh fyne-cross/dist/subman-macos-arm64.zip
ls -lh fyne-cross/dist/subman-linux-amd64.tar.xz
ls -lh fyne-cross/dist/subman-linux-arm64.tar.xz
ls -lh fyne-cross/dist/subman-windows-amd64.zip
ls -lh fyne-cross/dist/subman-windows-arm64.zip

echo ""
echo "📝 Creating GitHub release $VERSION..."
gh release create "$VERSION" --title "$VERSION" --generate-notes

echo ""
echo "⬆️  Uploading binaries to GitHub..."
gh release upload "$VERSION" \
  fyne-cross/dist/subman-macos-amd64.zip \
  fyne-cross/dist/subman-macos-arm64.zip \
  fyne-cross/dist/subman-linux-amd64.tar.xz \
  fyne-cross/dist/subman-linux-arm64.tar.xz \
  fyne-cross/dist/subman-windows-amd64.zip \
  fyne-cross/dist/subman-windows-arm64.zip

echo ""
echo "🎉 Release $VERSION published successfully!"
echo "🔗 https://github.com/douglasbarnum-cmyk/subman/releases/tag/$VERSION"
