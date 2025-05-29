Zadanie 1

Piotr Pepa

Punkt 1

Stworzona zostala aplikacja internetowa w jezyku GO (jezyk wybrany ze wzgledu na latwosc produkcji i maly rozmiar pliku + latwe linkowanie statyczne).
Aplikacja sklada sie z pliku main.go oraz index.html (dostpne na github)

----------------------------------------------------------------------------------------------------------------------------------------

Punkt 2

# Etap 1: Budowanie i kompresja binarki

# pobranie obrazu go z alpine
FROM golang:1.21-alpine AS builder
# katalog roboczy wewnatrz kontenera
WORKDIR /app
#przekopiowanie plikow z komputera do katalogu roboczego w kontenerze
COPY go.mod .
COPY main.go .
COPY index.html .
# zainstalowanie programu do kompresji binarek, zainstalowanie certyfikatow potrzebnych do HTTPS, wylanczenie zaleznosci systemowych (binarka bedzie statyczna), usuwanie debug info, plik wynikowy będzie nazywal sie main, skompresowanie binarki najwolniejsze ale dajace najmniejszy rozmiar
RUN apk add --no-cache upx ca-certificates && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main . && \
    upx --best --lzma main

# Etap 2: Distroless (zamiast scratch)

# distroless (bez bash, apk,  /bin, zawiera tylko najpotrzebniejsze katalogi i certyfikaty)
FROM gcr.io/distroless/static:nonroot

# Etykiety w stylu OCI (dalem wiecej niz wymagane minimum)
LABEL org.opencontainers.image.authors="Piotr Pepa <peter@trudne.eu>"
LABEL org.opencontainers.image.title="Pogodynka"
LABEL org.opencontainers.image.description="Aplikacja pogodowa w Go z danymi z dobrapogoda24.pl"
LABEL org.opencontainers.image.version="1.0"

# katalog roboczy wewnatrz kontenera
WORKDIR /app

# kopiuje skompilowana i skompresowana binarke oraz plik html z etapu pierwszego
COPY --from=builder /app/main /app/
COPY --from=builder /app/index.html /app/

# zadeklarowanie uzycia portu 8080
EXPOSE 8080

# healthcheck (czeka 5 s od startu, co 30 sekund, na maks 5 sekund, jesli 3 razy test nie przejdzie to jest uznane za unhealthy)
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD ["/app/main", "-healthcheck"]

# uruchamia main po uruchomieniu kontenera
ENTRYPOINT ["/app/main"]

-------------------------------------------------------------------------------------------------------------------------------------------------------------------

Punkt 3

zbudowanie obrazu:
docker build -t go-apka-1 .

uruchomienie kontenera na podstawie zbudowanego obrazu:
docker run -d --env-file .env -p 8080:8080 --name pogoda go-apka-1

sposobu uzyskania informacji z logow:
docker logs pogoda

sprawdzenia, ile warstw posiada zbudowany obraz:
docker image history go-apka-1

oraz jaki jest rozmiar obrazu:
docker image inspect go-apka-1


