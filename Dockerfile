FROM golang:1.22.2-alpine

# Install dependencies
RUN apk add --no-cache wkhtmltopdf ttf-dejavu fontconfig libgcc libstdc++ ca-certificates

WORKDIR /app

# Copy semua files
COPY . .

# Build aplikasi langsung
RUN go build -o app ./app/main.go

# Create symlink jika diperlukan
RUN ln -s /app/internal /internal 

# Debug permissions
RUN ls -la && chmod +x app && ls -la app

EXPOSE 80

# Gunakan exec form
CMD ["/app/app"]