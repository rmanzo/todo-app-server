run:
	go run ./cmd/web

buildserverlinux:
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/web

buildservermac:
	GOOS=darwin GOARCH=amd64 go build -o bin/server-mac ./cmd/web