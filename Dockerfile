# Menggunakan image Golang berbasis Debian untuk kompatibilitas lebih baik
FROM golang:1.22.2

# Environment variables untuk CompileDaemon
ENV GO111MODULE=on \
    CGO_ENABLED=0

# Menentukan working directory di dalam container
WORKDIR /app

# Membuat folder untuk hasil build
RUN mkdir -p /build

# Menyalin semua file dari project ke dalam container
COPY . .

# Menyalin go.mod dan go.sum untuk dependency management
COPY go.mod go.sum ./

# Mengunduh semua dependency
RUN go mod download

# Menginstal CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

# Mengekspos port aplikasi
EXPOSE 80

# Menjalankan CompileDaemon agar otomatis rebuild dan restart saat ada perubahan kode
ENTRYPOINT ["CompileDaemon", "-build=go build -o /build/app ./app/main.go", "-command=/build/app"]
