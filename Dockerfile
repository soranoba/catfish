FROM golang:1.19 AS builder

RUN apt-get update && \
    apt-get install -y npm && \
    apt-get purge -y

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
ARG GOFLAGS
RUN GOFLAGS="${GOFLAGS}" make release

#==========================

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/catfish /bin/catfish
COPY --from=builder /app/bin/config.yml /etc/catfish/config.yml

ENTRYPOINT ["catfish", "--config", "/etc/catfish/config.yml"]
