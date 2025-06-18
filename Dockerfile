# Stage 1: Install wkhtmltopdf di Alpine 3.14
FROM alpine:3.14 AS wkhtml

RUN apk add --no-cache \
    wkhtmltopdf \
    ttf-dejavu \
    fontconfig

# Stage 2: Build application
FROM golang:1.22.2-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

# Copy go.mod dan go.sum terlebih dahulu untuk caching layer
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Buat symlink jika diperlukan untuk build
RUN ln -s /build/internal /internal

# Build aplikasi untuk production
RUN go build -ldflags="-s -w" -o app ./app/main.go

# Stage 3: Final lightweight image
FROM alpine:3.14

WORKDIR /app

# Install dependencies yang diperlukan untuk runtime
RUN apk add --no-cache \
    libgcc \
    libstdc++ \
    ttf-dejavu \
    fontconfig \
    ca-certificates

# Copy wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# Copy compiled binary dari builder stage
COPY --from=builder /build/app /app/app

# Verify file exists and set permissions
RUN ls -la && chmod +x app

# Buat symlink jika diperlukan untuk runtime
RUN mkdir -p /internal && ln -s /app/internal /internal

EXPOSE 80

# Use non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Run the application
CMD ["/app/app"]