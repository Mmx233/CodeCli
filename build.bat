@echo off
SET /p v=Version:
release ./cmd/code --os linux,windows,darwin --arch amd64,386,arm64 --ldflags="-X main.Version=%v%"