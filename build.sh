${HOME}/go/bin/gox -osarch="linux/amd64 windows/amd64 darwin/amd64 darwin/arm64" -output="swed_{{.OS}}_{{.Arch}}" -ldflags="-s -w" ./app/swed