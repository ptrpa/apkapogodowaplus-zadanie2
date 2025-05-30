```yaml
//nazwa która pojawi się w github actions jako nazwa workflow
name: Build & Push Multiarch Docker Image with Security Scan

//przypadki kiedy to się uruchamia
on:
//push na gałąź main
  push:
    branches: [ main ]
//i sam pull_request na gałąź main
  pull_request:
    branches: [ main ]

//zadania które są do wykonania (i które się pokazują na diagramie), w tym przypadku są dwa
jobs:
  build-and-scan:
    name: Build and Scan
    runs-on: ubuntu-22.04

    steps:
      //pobranie kodu źródłowego z gałęzi
      - name: Checkout repository
        uses: actions/checkout@v4
      //ustawienie qemu potrzebnego do emulacji architektur
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      //stworzenie instancji buildx potrzebnej do m.in. multiarch oraz do cache
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      //zalogowanie do dockerhub aby móc tam wrzucać cache, użyto sekretów
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      //zalogowanie do github container registry aby wrzucać tam ukończony obraz, także użyto skrtów
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ secrets.GHCR_USERNAME }}
          password: ${{ secrets.GHCR_TOKEN }}
      //zbudowanie obrazu tylko w celu skanowania trivy, wystarczy dla amd, nie pushuje, natomiast taguje obraz, korzysta z cache z dockerhuba (oczywiście o ile wcześniej już zbudowano cokolwiek)
      - name: Build image for scanning
        uses: docker/build-push-action@v5
        with:
          context: .
          push: false
          load: true
          platforms: linux/amd64
          tags: |
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:latest
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:v1.0.0
          cache-from: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache
          cache-to: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache,mode=max
      //wykonanie sprawdzenia obrazu przez trivy, przejdzie do kolejnego joba jedynie jeśli nie ma poważnych zagrożeń, wynik zapisuje się też do pliku txt
      - name: Run Trivy security scan
        uses: aquasecurity/trivy-action@0.11.2
        with:
          image-ref: ghcr.io/ptrpa/apkapogodowaplus-zadanie2:latest
          format: table
          output: trivy-report.txt
          github-pat: ${{ secrets.GHCR_TOKEN }}
          exit-code: 1
          severity: HIGH,CRITICAL
        //przesłanie wyniku skanu jako artefakt niezależnie czy skan przejdzie pozytywnie czy nie
      - name: Upload Trivy report artifact
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: trivy-report
          path: trivy-report.txt
  //drugi job uruchomi się tylko po sukcesie pierwszego
  push:
    name: Push Multiarch Image
    runs-on: ubuntu-22.04
    needs: build-and-scan
//dodatkowe zabezpieczenie żeby się nie uruchomił jak wcześniej coś było nie tak
    if: success()
//pięć kroków analogicznie do pierwszego joba
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ secrets.GHCR_USERNAME }}
          password: ${{ secrets.GHCR_TOKEN }}

      - name: Build and push multiarch image
        uses: docker/build-push-action@v5
        with:
          context: .
//w przeciwieństwie do poprzedniego joba obraz po zbudowaniu będzie opublikowany na ghcr
          push: true
//tutaj już obraz na wiele architektur
          platforms: linux/amd64,linux/arm64
//otagowanie obrazu według https://semver.org/ (było na wykładzie "z przemysłu" na szkieletach programistycznych w aplikacjach internetowych)
          tags: |
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:latest
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:v1.0.0
//korzystanie z cache z dockerhub
          cache-from: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache
          cache-to: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache,mode=max

```
