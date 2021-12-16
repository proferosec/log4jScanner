PROJECT_NAME=log4jScanner
GOPATH=$(shell go env GOPATH)

VERSION=$(shell $(GOPATH)/bin/semver get release)
BUILD_TIME=$(shell TZ=UTC date -u '+%Y-%m-%d_%I:%M:%S%p')

all: clean init build release

build: build-windows build-darwin build-linux
release: release-windows release-darwin release-linux

test:
	go test .

init:
	go get -u github.com/maykonlf/semver-cli/cmd/semver 

upver:
	$(GOPATH)/bin/semver up alpha

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "build/windows/$(PROJECT_NAME)-$(VERSION).exe" -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o "build/darwin/$(PROJECT_NAME)-$(VERSION)" -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "build/linux/$(PROJECT_NAME)-$(VERSION)" -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

release-windows: release-dir
	zip -j release/$(PROJECT_NAME)-windows.zip build/windows/$(PROJECT_NAME)-$(VERSION).exe
	echo $(shell shasum -a 256 build/windows/$(PROJECT_NAME)-$(VERSION).exe | cut -f1 -d" ") $(PROJECT_NAME)-$(VERSION).exe > release/windows.sha256.txt

release-darwin: release-dir
	zip -j release/$(PROJECT_NAME)-darwin.zip build/darwin/$(PROJECT_NAME)-$(VERSION)
	echo $(shell shasum -a 256 build/darwin/$(PROJECT_NAME)-$(VERSION) | cut -f1 -d" ") $(PROJECT_NAME)-$(VERSION) > release/darwin.sha256.txt

release-linux: release-dir
	zip -j release/$(PROJECT_NAME)-linux.zip build/linux/$(PROJECT_NAME)-$(VERSION)
	echo $(shell shasum -a 256 build/linux/$(PROJECT_NAME)-$(VERSION) | cut -f1 -d" ") $(PROJECT_NAME)-$(VERSION) > release/linux.sha256.txt

release-dir:
	mkdir release || true

clean:
	rm -rf ./build || true
	rm -rf ./release || true
	rm *.log || true
	rm *.csv || true
