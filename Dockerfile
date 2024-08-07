FROM golang:1.22-alpine as builder

WORKDIR /app
COPY  . .

RUN GOOS=linux go build -o /app/go-echo ./cmd/go-echo/*.go

FROM ubuntu:22.04

ARG image_source
LABEL org.opencontainers.image.source ${image_source}

COPY --from=builder /app/go-echo /usr/local/bin/

CMD ["go-echo"]
