FROM golang:1.17 AS builder
WORKDIR /go/src
COPY . .
RUN make build

FROM alpine:3.14
ARG APP
WORKDIR /app
COPY --from=builder /go/src/tmp/${APP} .
COPY --from=builder /go/src/config.toml .
CMD ["./${APP}", "server"]
