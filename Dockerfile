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

# PERBAIKAN 1: Build langsung ke binary executable
RUN go build -o main ./app/main.go  # Ubah output nama

# Stage 3: Final image
FROM alpine:3.14

# Install runtime dependencies
RUN apk add --no-cache libgcc libstdc++ ttf-dejavu fontconfig ca-certificates

# Copy wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# PERBAIKAN 2: Buat direktori khusus untuk app
WORKDIR /app
RUN mkdir -p /app/bin

# Copy aplikasi dari builder
COPY --from=builder /app/main /app/bin/main  # Path konsisten

# PERBAIKAN 3: Set permission secara eksplisit
RUN chmod -R 755 /app/bin

# PERBAIKAN 4: Gunakan user non-root
RUN adduser -D myuser
USER myuser

EXPOSE 80

# PERBAIKAN 5: Gunakan full path ke binary
CMD ["/app/bin/main"]