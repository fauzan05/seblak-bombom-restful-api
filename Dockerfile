# Stage 1: Install wkhtmltopdf
FROM alpine:3.14 AS wkhtml
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig

# Stage 2: Build app
FROM golang:1.22.2-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

# Copy dan build project
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build aplikasi ke dalam folder /build
RUN go build -o app ./app/main.go && ls -la

# Stage 3: Final image
FROM alpine:3.14

# Install runtime dependencies
RUN apk add --no-cache libgcc libstdc++ ttf-dejavu fontconfig ca-certificates

# Copy wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# Create workdir
WORKDIR /app

# Copy aplikasi dari builder
COPY --from=builder /build/app .

# Debug: check file exists
RUN ls -la /app/

# Set permission
RUN chmod +x app

EXPOSE 80

CMD ["./app"]