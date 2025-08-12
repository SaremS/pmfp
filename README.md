# Poor man's forward proxy
[![Go Tests](https://github.com/SaremS/pmfp/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/SaremS/pmfp/actions/workflows/go.yml)

## Summary
(currently http only)

Forward proxy with minimal features (currently only adds custom headers to each request)
## Usage
Currenty
```
go run cmd/main.go --config {config-yaml; see examples} --port {target port to bind proxy to}
```
