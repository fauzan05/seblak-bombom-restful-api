FROM surnet/alpine-wkhtmltopdf:3.8-0.12.5-full as wkhtmlbuilder

FROM golang:1.22-alpine

# Install needed packages
RUN apk update && apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu

# Copy wkhtmltopdf
COPY --from=wkhtmlbuilder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmlbuilder /bin/wkhtmltoimage /bin/wkhtmltoimage

WORKDIR /app

# Copy seluruh source code
COPY . .

# Debug: lihat struktur folder
RUN ls -la && ls -la app/

# Build aplikasi dari app/main.go ke root folder dengan nama 'api'
RUN go build -o api ./app/main.go

# Debug: verifikasi api telah di-build
RUN ls -la && chmod +x api

# Create symlink jika diperlukan
RUN ln -s /app/internal /internal || true

EXPOSE 80

# Gunakan executable baru
CMD ["/app/api"]