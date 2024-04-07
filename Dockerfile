FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /seblak-bombom

COPY . .

RUN go mod tidy

RUN go build -o /seblak-bombom/seblak_bombom ./app/main.go

FROM alpine:3.19

WORKDIR /seblak-bombom

RUN apk update && \
    apk add --no-cache mariadb-client

COPY --from=builder /seblak-bombom/seblak_bombom ./
COPY config.json ./
EXPOSE 8000
ENV database_url=mysql://root@tcp(localhost:3306)/seblak_bombom

ENTRYPOINT [  "./seblak_bombom" ]