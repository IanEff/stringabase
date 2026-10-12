# stringabase

A barebones, concurrency-safe, in-memory string key-value store in Go, served over HTTP.

Structure loosely taken from "Cloud Native Go" by Matthew A. Titmus and "The Deeper Love of Go" by John Arundel, with a stab at commenting and design from the approach in John Ousterhout's "A Philosophy of Software Design."

## Usage
1. Start the server in one terminal:
```bash
go run ./cmd/stringabase
```

2. In another terminal, store a value:
```bash
curl -i 'http://localhost:4000/set?somekey=somevalue'
```
You'll get an HTTP 200 status code and an empty response body.

3. Retrieve your value:
```bash
curl -w '\n' 'http://localhost:4000/get?key=somekey'
```
You'll get your value back.

4. Stop the server by hitting `ctrl-c` in its terminal.

## Layout

I've organized the Stringabase into three packages, to keep concerns separate:

```
.
├── cmd/
│   └── stringabase/    the composition root, where all defaults are baked in
└── internal/
    ├── server/         the HTTP server and endpoints
    └── stringabase/    the key-value store and its accessors
```
