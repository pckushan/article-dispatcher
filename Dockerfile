FROM golang:1.18-alpine3.15 AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go build -o article-dispatcher *.go

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /src
COPY --from=builder /app /src
ENV HTTP_SERVER_HOST=8888
ENV LOG_LEVEL="TRACE"

EXPOSE $HTTP_SERVER_PORT
EXPOSE $METRICS_PORT

CMD ["./article-dispatcher"]
