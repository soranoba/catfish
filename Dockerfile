FROM golang:1.19 AS builder

WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
RUN make release

#==========================

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/catfish /bin/catfish
COPY --from=builder /app/bin/config.yml /etc/cartfish/config.yml

ENTRYPOINT ["catfish", "--config", "/etc/cartfish/config.yml"]
