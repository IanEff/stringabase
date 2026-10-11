# stringabase

[![Go Version](https://img.shields.io/github/go-mod/go-version/IanEff/stringabase?logo=go)](https://go.dev)

A barebones, concurency-safe database server in go.

Structure loosely taken from "Cloud Native Go", by Matthew A. Titmus, and "For the Deeper Love of Go" by John Arundel, with a stab at commenting and design from the approach in John Ousterhout's "A Philsophy of Software Design."

## Usage
1. Start the server on in one terminal: 
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
curl -w 'n' 'http://localhost:4000/get?key=somekey'
```
You'll get another HTTP 200 and your value back.
4. Stop the server by hitting `ctrl-c` in its terminal.
