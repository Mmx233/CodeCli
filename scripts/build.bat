@echo off
SET /p v=Version:
release ./cmd/code --os windows --arch amd64,386,arm64 --ldflags="-X main.Version=%v%" -c tar.gz