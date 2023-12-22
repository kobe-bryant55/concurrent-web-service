# Builder
FROM --platform=linux/amd64/v8 golang:1.21.4-alpine3.17 as builder

RUN apk add alpine-sdk

ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOARCH="amd64"
ENV GOOS=linux

WORKDIR /build

COPY go.* ./

RUN go mod download

COPY . ./

RUN go build -tags musl -o concurrenct-web-service ./cmd/api

# Application container
FROM --platform=linux/amd64/v8 alpine:3.17.2

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/concurrenct-web-service /app/concurrenct-web-service

CMD ["/app/concurrenct-web-service"]