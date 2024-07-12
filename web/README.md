# The kube-score website

This directory contains the source for kube-score.com.

kube-score is compiled to WASM (WebAssembly) for usage in browsers, and kube-score.com is running purely in the browser.

## Building

```bash
# Compile the wasm module (run from repo root)
GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./cmd/wasm/main.go
```

```bash
# Start a webserver serving files
python3 -m http.server
```

## Hosting

The site is built and hosted on Vercel. See vercel.json for build configuration.

