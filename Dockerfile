FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine AS runner

RUN apk update && \
    apk add --no-cache tzdata ca-certificates && \
    rm -rf /var/cache/apk/*

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /usr/local/src/bin/app /app

COPY .env /.env
COPY config/config.yaml /config/config.yaml
COPY files /files

EXPOSE 8080

CMD ["/app"]