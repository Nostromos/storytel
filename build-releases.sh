#!/bin/bash

# Clean and create release directories
rm -rf releases/*
mkdir -p releases/macos-arm64
mkdir -p releases/macos-amd64
mkdir -p releases/linux-amd64
mkdir -p releases/linux-arm64
mkdir -p releases/windows-amd64

# Build for each platform
GOOS=darwin GOARCH=arm64 go build -o releases/macos-arm64/storytel ./cmd/storytel
GOOS=darwin GOARCH=amd64 go build -o releases/macos-amd64/storytel ./cmd/storytel
GOOS=linux GOARCH=amd64 go build -o releases/linux-amd64/storytel ./cmd/storytel
GOOS=linux GOARCH=arm64 go build -o releases/linux-arm64/storytel ./cmd/storytel
GOOS=windows GOARCH=amd64 go build -o releases/windows-amd64/storytel.exe ./cmd/storytel

echo "Build complete! Binaries are in the releases/ directory"