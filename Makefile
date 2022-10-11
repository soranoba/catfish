VERSION     := $(shell git describe --always --tags --abbrev=10)
CMDS         = ls -d cmd/* | xargs -I@ basename @
CONFIG_PKG   = github.com/soranoba/catfish/pkg/config

build: build-js build-app

build-app:
	${CMDS} | xargs -I@ go build -ldflags "-X ${CONFIG_PKG}.AppVersion=${VERSION}" -o bin/@ ./cmd/@

build-js:
	cd cmd/catfish/static && npm ci && npm run build

release-app: build-js
	${CMDS} | xargs -I@ go build -ldflags "-s -w -X ${CONFIG_PKG}.AppVersion=${VERSION}" -o bin/@ -a ./cmd/@

start:
	cd cmd/catfish/static && npm ci && npm run build
	go run -a ./cmd/catfish -- --config=./bin/config.yml

test:
	go test ./... -count=1

format:
	gofmt -w ./

lint:
	gofmt -d ./
	test -z $(shell gofmt -l ./)
