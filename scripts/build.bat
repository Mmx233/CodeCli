@echo off
SET /p v=Version:
release ./cmd/code --os windows,linux --extra-arches --ldflags="-X main.Version=%v%" -c tar.gz