VERSION     := $(shell git describe --always --tags --abbrev=10)
CMDS         = ls -d cmd/* | xargs -I@ basename @
CONFIG_PKG   = github.com/soranoba/catfish/pkg/config
BUILDVCS    := $(shell echo $${BUILDVCS:-true})

build:
	cd cmd/catfish/static && npm ci && npm run build
	${CMDS} | xargs -I@ go build -ldflags "-X ${CONFIG_PKG}.AppVersion=${VERSION}" -buildvcs=${BUILDVCS} -o bin/@ ./cmd/@

release:
	cd cmd/catfish/static && npm ci && npm run build
	${CMDS} | CGO_ENABLED=0 GOOS=linux GOARCH=amd64 xargs -I@ \
		go build -ldflags "-s -w -X ${CONFIG_PKG}.AppVersion=${VERSION}" -buildvcs=${BUILDVCS} -o bin/@ -a ./cmd/@

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
