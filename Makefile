CMDS = ls -d cmd/* | xargs -I@ basename @

build:
	${CMDS} | xargs -I@ go build -o bin/@ ./cmd/@

start:
	go run -a ./cmd/catfish

test:
	go test ./... -count=1

format:
	gofmt -w ./

lint:
	gofmt -d ./
	test -z $(shell gofmt -l ./)
