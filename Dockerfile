FROM node:16 AS js-builder

WORKDIR /app/cmd/catfish/static
COPY ./cmd/catfish/static/package.json .
COPY ./cmd/catfish/static/package-lock.json .
RUN npm ci
COPY ./cmd/catfish/static/ .
WORKDIR /app
COPY ./Makefile ./
RUN make build-js

#==========================

FROM golang:1.19 AS app-builder

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
COPY --from=js-builder /app/cmd/catfish/static/public ./cmd/catfish/static/
ARG GOFLAGS
ARG GOOS=linux
ARG GOARCH=amd64
RUN GOFLAGS="${GOFLAGS}" GOOS="${GOOS}" GOARCH="${GOARCH}" make release-app
#==========================

FROM alpine:latest

WORKDIR /app
COPY --from=app-builder /app/bin/catfish /bin/catfish
COPY --from=app-builder /app/bin/config.yml /etc/catfish/config.yml

ENTRYPOINT ["catfish", "--config", "/etc/catfish/config.yml"]
