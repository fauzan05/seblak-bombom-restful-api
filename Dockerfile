FROM golang:1.22.2-alpine

# Environment variables which CompileDaemon requires to run
ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

RUN mkdir "/build"

COPY . .

COPY go.mod .

# bisa juga menggunakan tidy
RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 80

ENTRYPOINT CompileDaemon -build="go build -o /build/app ./app/main.go" -command="/build/app"