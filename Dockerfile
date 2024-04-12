FROM golang:1.22.2-alpine

# Environment variables which CompileDaemon requires to run
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

RUN mkdir "/build"

COPY . .

COPY go.mod .

RUN go mod download
# bisa juga menggunakan tidy

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 8000

ENTRYPOINT CompileDaemon -directory="/app"  -build="go build -o /build/app ./app/main.go" -command="/build/app"