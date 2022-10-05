FROM node:16 AS js-builder

WORKDIR /app/cmd/catfish/static
COPY ./cmd/catfish/static/package.json .
COPY ./cmd/catfish/static/package-lock.json .
RUN npm ci
COPY ./cmd/catfish/static/ .
WORKDIR /app
COPY ./Makefile ./
RUN make release-js

#==========================

FROM golang:1.19 AS builder

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
COPY --from=js-builder /app/cmd/catfish/static/public ./cmd/catfish/static/
ARG GOFLAGS
ARG GOARCH=amd64
RUN GOFLAGS="${GOFLAGS}" GOARCH="${GOARCH}" make release-app
#==========================

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/catfish /bin/catfish
COPY --from=builder /app/bin/config.yml /etc/catfish/config.yml

ENTRYPOINT ["catfish", "--config", "/etc/catfish/config.yml"]
