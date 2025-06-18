# Stage 1: Install wkhtmltopdf
FROM alpine:3.14 AS wkhtml
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig

# Stage 2: Build app
FROM golang:1.22.2-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -o app ./app/main.go

# Stage 3: Final image - SANGAT SEDERHANA
FROM alpine:3.14

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache libgcc libstdc++ ttf-dejavu fontconfig ca-certificates

# Copy wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# Copy compiled app
COPY --from=builder /build/app .

# PENTING: Set executable permission
RUN chmod 777 /app/app && ls -la

# Create symlink jika diperlukan
RUN ln -s /app/internal /internal

EXPOSE 80

# SANGAT PENTING: Gunakan root user
# USER root

CMD ["/app/app"]