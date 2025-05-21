FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . ./

COPY config/config.yaml /tmp/config.yaml

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine AS runner

RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /config

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /tmp/config.yaml /config/config.yaml
COPY .env /.env
COPY files /files

EXPOSE 1111

CMD ["/app"]