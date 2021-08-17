GOTOOLS =	github.com/mitchellh/gox \
			github.com/Masterminds/glide \
			github.com/rigelrozanski/shelldown/cmd/shelldown
INCLUDE = -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf

all: install

build: gen_version
	go build ./cmd/...

# Build binaries for Linux platform.
linux: gen_version
	integration/docker/build/build.sh force

windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o build/windows/edgecore.exe ./cmd/theta-eth-rpc-adaptor

docker: 
	integration/docker/node/build.sh force

install: gen_version release

release:
	go install ./cmd/...

debug:
	go install -race ./cmd/...

clean:
	@rm -rf ./build

gen_doc:
	cd ./docs/commands/;go build -o generator.exe; ./generator.exe

BUILD_DATE := `date -u`
GIT_HASH := `git rev-parse HEAD`
VERSION_NUMER := `cat version/version_number.txt`
VERSIONFILE := version/version_generated.go

gen_version:
	@echo "package version" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  Timestamp = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo "  Version = \"$(VERSION_NUMER)\"" >> $(VERSIONFILE)
	@echo "  GitHash = \"$(GIT_HASH)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)

.PHONY: all build install clean
