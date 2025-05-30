```yaml
name: Build & Push Multiarch Docker Image with Security Scan

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build-and-scan:
    name: Build and Scan
    runs-on: ubuntu-22.04

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

      - name: Run Trivy security scan
        uses: aquasecurity/trivy-action@0.11.2
        with:
          image-ref: ghcr.io/ptrpa/apkapogodowaplus-zadanie2:latest
          format: table
          output: trivy-report.txt
          github-pat: ${{ secrets.GHCR_TOKEN }}
          exit-code: 1
          severity: HIGH,CRITICAL

      - name: Upload Trivy report artifact
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: trivy-report
          path: trivy-report.txt

  push:
    name: Push Multiarch Image
    runs-on: ubuntu-22.04
    needs: build-and-scan
    if: success()

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
          push: true
          platforms: linux/amd64,linux/arm64
          tags: |
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:latest
            ghcr.io/ptrpa/apkapogodowaplus-zadanie2:v1.0.0
          cache-from: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache
          cache-to: type=registry,ref=s99656/apkapogodowaplus-zadanie2:buildcache,mode=max

```
