# Stage 1: Build binary
FROM golang:1.22.2-alpine3.14 AS builder

WORKDIR /app

# Aktifkan CGO kalau wkhtmltopdf membutuhkan fontconfig & sejenisnya
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache \
    git \
    libc6-compat \
    build-base \
    wkhtmltopdf \
    ttf-dejavu \
    fontconfig

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary ke /build/app
RUN mkdir /build && \
    go build -o /build/app ./app/main.go

# ─────────────────────────────────────────────────────────────────────

# Stage 2: Minimal runtime container
FROM alpine:3.14

WORKDIR /app

# Copy binary dan wkhtmltopdf dari builder
COPY --from=builder /build/app /app/app
COPY --from=builder /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=builder /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=builder /usr/share/fonts /usr/share/fonts
COPY --from=builder /etc/fonts /etc/fonts

# Runtime deps
RUN apk add --no-cache \
    libgcc \
    fontconfig \
    ttf-dejavu

EXPOSE 80

CMD ["./app"]
