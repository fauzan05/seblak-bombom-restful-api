FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /seblak-bombom

COPY . .
# bisa juga menggunakan tidy
RUN go mod download

RUN go build -o /seblak-bombom/seblak_bombom ./app/main.go

FROM alpine:3.19

WORKDIR /seblak-bombom

# RUN apk update && \
#     apk add --no-cache go 

COPY --from=builder /seblak-bombom/seblak_bombom ./
COPY config.json ./

# COPY database/ ./
# RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# RUN ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

EXPOSE 8000

ENTRYPOINT [ "./seblak_bombom" ]