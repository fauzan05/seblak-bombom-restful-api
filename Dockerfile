# Stage 1: Install wkhtmltopdf di Alpine 3.14
FROM alpine:3.14 AS wkhtml

RUN apk add --no-cache \
    wkhtmltopdf \
    ttf-dejavu \
    fontconfig

# Stage 2: Gunakan image Go dengan Alpine terbaru
FROM golang:1.22.2-alpine

# ─── Env untuk CompileDaemon ───────────────────────────────────────
ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

# ─── Salin wkhtmltopdf dari Stage 1 ────────────────────────────────
COPY --from=wkhtml /usr/bin/wkhtmltopdf /usr/bin/wkhtmltopdf
COPY --from=wkhtml /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=wkhtml /usr/share/fonts /usr/share/fonts

# ─── Install Dependensi Tambahan ───────────────────────────────────
RUN apk add --no-cache \
    libgcc \
    ttf-dejavu \
    fontconfig

# ─── Copy Project & Buat Symlink untuk Internal Directory ──────────
COPY . .

# Buat symlink agar "../internal" di kode mengarah ke "/app/internal"
RUN ln -s /app/internal /internal

# ─── Build tools & app ─────────────────────────────────────────────
RUN mkdir /build

COPY go.mod .
RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest

EXPOSE 80

CMD ["./app"]
