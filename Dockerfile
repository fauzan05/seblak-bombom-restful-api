FROM surnet/alpine-wkhtmltopdf:3.8-0.12.5-full as wkhtmlbuilder

FROM golang:1.22-alpine

RUN apk update && apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu

COPY --from=wkhtmlbuilder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmlbuilder /bin/wkhtmltoimage /bin/wkhtmltoimage

WORKDIR /app

# Copy semua files
COPY . .

# Debug
RUN ls -la

# Create wrapper script
RUN echo '#!/bin/sh' > start.sh && \
    echo 'cd /app && ./main' >> start.sh && \
    chmod +x start.sh

# Create symlink jika diperlukan
RUN ln -s /app/internal /internal || true

EXPOSE 80

# Gunakan script wrapper
CMD ["/app/start.sh"]