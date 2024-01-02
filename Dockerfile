# Builder
FROM golang:1.20.12-alpine3.19 as builder

RUN apk update && apk upgrade && \
   apk --update add git make bash build-base

WORKDIR /app
ENV CGO_ENABLED=0

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
   apk --update --no-cache add tzdata && \
   mkdir /apps

WORKDIR /app

EXPOSE 8000

COPY --from=builder /app/target/github-webhook-telegram-bot /app
COPY --from=builder /app/.env /app

WORKDIR /app


CMD /app/github-webhook-telegram-bot