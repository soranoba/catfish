VERSION     := $(shell git describe --always --tags --abbrev=10)
CMDS         = ls -d cmd/* | xargs -I@ basename @
CONFIG_PKG   = github.com/soranoba/catfish-server

build:
	${CMDS} | xargs -I@ go build -ldflags "-X ${CONFIG_PKG}.AppVersion=${VERSION}" -o bin/@ ./cmd/@

start:
	go run -a ./cmd/catfish -- --config=./bin/config.yml

test:
	go test ./... -count=1

format:
	gofmt -w ./

lint:
	gofmt -d ./
	test -z $(shell gofmt -l ./)
