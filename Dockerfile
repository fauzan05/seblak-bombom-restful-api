# builder image untuk wkhtmltopdf
FROM surnet/alpine-wkhtmltopdf:3.8-0.12.5-full as wkhtmlbuilder

# Build golang app
FROM golang:1.22-alpine AS builder

WORKDIR /go/src/app

# Copy go.mod dan go.sum terlebih dahulu (jika ada)
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY . .

# Build aplikasi
RUN go build -o app ./app/main.go

# Final image
FROM golang:1.11-alpine3.8

# Install needed packages
RUN  echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/main" > /etc/apk/repositories \
     && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/community" >> /etc/apk/repositories \
     && apk update && apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      libssl1.0 \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu \
      ttf-droid \
      ttf-freefont \
      ttf-liberation \
      ttf-ubuntu-font-family \
    && apk add --no-cache --virtual .build-deps \
      msttcorefonts-installer \
    \
    # Install microsoft fonts
    && update-ms-fonts \
    && fc-cache -f \
    \
    # Clean up when done
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* \
    && apk del .build-deps

# Copy wkhtmltopdf
COPY --from=wkhtmlbuilder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmlbuilder /bin/wkhtmltoimage /bin/wkhtmltoimage

WORKDIR /go/src/app

# Copy built app dari builder stage
COPY --from=builder /go/src/app/app .

# Debug: cek app dan beri permission
RUN ls -la && chmod +x app

# Create symlink jika diperlukan
RUN ln -s /go/src/app/internal /internal || true

EXPOSE 80

# Gunakan absolute path
CMD ["/go/src/app/app"]