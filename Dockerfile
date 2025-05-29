# Etap 1: Budowanie i kompresja binarki
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .
COPY index.html .

RUN apk add --no-cache upx ca-certificates && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main . && \
    upx --best --lzma main

# Etap 2: Distroless (zamiast scratch)
FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.authors="Piotr Pepa <peter@trudne.eu>"
LABEL org.opencontainers.image.title="Pogodynka"
LABEL org.opencontainers.image.description="Aplikacja pogodowa w Go z danymi z dobrapogoda24.pl"
LABEL org.opencontainers.image.version="1.0"

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/index.html /app/

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD ["/app/main", "-healthcheck"]

ENTRYPOINT ["/app/main"]
