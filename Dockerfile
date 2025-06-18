# Stage 1: Install wkhtmltopdf
FROM alpine:3.14 AS wkhtml
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig

# Stage 2: Golang build
FROM golang:1.22.2-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

# wkhtmltopdf
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

RUN apk add --no-cache libgcc ttf-dejavu fontconfig

# Project code
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ln -s /app/internal /internal

# ✅ BUILD APP + beri permission execute
RUN go build -o /app/app ./app/main.go && chmod +x /app/app

EXPOSE 80

# ✅ Jalankan binary yang sudah bisa dieksekusi
CMD ["./app"]
