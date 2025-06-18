# Stage 1: Install wkhtmltopdf
FROM alpine:3.14 AS wkhtml
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig

# Stage 2: Build app
FROM golang:1.22.2-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

# Install wkhtmltopdf dari stage 1
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# Install dependencies
RUN apk add --no-cache libgcc ttf-dejavu fontconfig

# Copy dan build project
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN ln -s /app/internal /internal

# ✅ Build ke ./app dan beri izin eksekusi
RUN go build -o app ./app/main.go && chmod +x app

EXPOSE 80

CMD ["./app"]
