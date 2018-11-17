test:
	go get ./...
	go test -v ./...

build: test
	go get github.com/karalabe/xgo
	GOOS=linux GOARCH=amd64 go build -o bens-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o bens-macos-amd64 .
	GOOS=windows GOARCH=amd64 go get ./...
	GOOS=windows GOARCH=amd64 go build -o bens-windows-amd64.exe .

clean:
	rm -rf bens-*

.PHONY: build clean test
