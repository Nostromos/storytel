#!/bin/bash

VERSION=${1:-v1.0}

# First build the releases
echo "Building releases..."
./build-releases.sh

# Create archives directory
mkdir -p archives

# Package each platform
echo "Packaging releases..."

# macOS ARM64
tar -czf archives/storytel-${VERSION}-macos-arm64.tar.gz -C releases/macos-arm64 storytel

# macOS AMD64
tar -czf archives/storytel-${VERSION}-macos-amd64.tar.gz -C releases/macos-amd64 storytel

# Linux AMD64
tar -czf archives/storytel-${VERSION}-linux-amd64.tar.gz -C releases/linux-amd64 storytel

# Linux ARM64
tar -czf archives/storytel-${VERSION}-linux-arm64.tar.gz -C releases/linux-arm64 storytel

# Windows AMD64 (using zip)
cd releases/windows-amd64
zip ../../archives/storytel-${VERSION}-windows-amd64.zip storytel.exe
cd ../..

echo "Release archives created in archives/ directory:"
ls -la archives/