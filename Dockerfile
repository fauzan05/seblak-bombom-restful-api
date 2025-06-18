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

COPY . .

# Berikan permission pada file main yang sudah ada
RUN chmod +x main || true

EXPOSE 80

# Gunakan file main yang sudah ada
CMD ["/app/main"]