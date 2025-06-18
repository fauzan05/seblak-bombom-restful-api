# Stage 1: Install wkhtmltopdf
FROM alpine:3.14 AS wkhtml
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig

# Stage 2: Build app
FROM golang:1.22.2-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# PERBAIKAN: Pastikan output path konsisten
RUN go build -o /app/main ./app/main.go  # Ubah menjadi /app/main

# Stage 3: Final image
FROM alpine:3.14

# Install runtime dependencies
RUN apk add --no-cache libgcc libstdc++ ttf-dejavu fontconfig ca-certificates

# Copy wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# Buat direktori
RUN mkdir -p /app/bin

# PERBAIKAN: Copy dengan path yang benar
COPY --from=builder /app/main /app/bin/main

# Set permission
RUN chmod 755 /app/bin/main

# Gunakan user non-root
RUN adduser -D -g '' myuser
USER myuser

WORKDIR /app

EXPOSE 80

CMD ["/app/bin/main"]