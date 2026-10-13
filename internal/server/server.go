// Package server exposes [stringabase.Stringabase] as an HTTP API. The endpoints are documented in
// [NewMuxer].
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ianeff/stringabase/internal/stringabase"
)

// ListenAndServe serves sb out over HTTP at addr on endpoints described in [NewMuxer], until ctx is canceled.
// sb is an initialized [stringabase.Stringabase]. addr takes any form net/http accepts ("localhost:4000").
// When ctx is canceled, ListenAndServe drops in-flight requests and returns nil. If it cannot serve, it
// immediately returns a non-nil error.
func ListenAndServe(ctx context.Context, addr string, sb *stringabase.Stringabase) error {
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      NewMuxer(sb),
	}
	stop := context.AfterFunc(ctx, func() { srv.Close() })
	defer stop()
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) && ctx.Err() != nil {
		return nil
	}
	return err
}

// NewMuxer returns a handler for the HTTP API, backed by sb.
// It exposes the following endpoints:
//
// /get retrieves a value from the Stringabase by key. Returns a StatusOK and the value if the key is
// present, even if that value is an empty string.
// Returns a StatusNotFound if the key is not found.
//
//	/get?key=somekey
//
// /set stores a key-value string pair in the Stringabase. It is idempotent and will replace any existing
// value for a key.  One request can set several keys. If a key repeats, the first value wins.
// Returns a StatusOK with an empty body.
//
//	/set?somekey=somevalue&otherkey=someothervalue
//
// Keys and values must be URL-encoded.
func NewMuxer(sb *stringabase.Stringabase) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", handleGet(sb))
	mux.HandleFunc("/set", handleSet(sb))
	return mux
}

// handleGet returns a handler that retrieves the query parameter named "key" from the URI and attempts to
// retrieve its value from sb.
//
//	/get?key=somekey
//
// Returns an HTTP StatusOK and the value if the key "somekey" is present, even if that value is an empty
// string, and a StatusNotFound if the key is not found.
func handleGet(sb *stringabase.Stringabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		key := query.Get("key")
		value, ok := sb.Get(key)
		if !ok {
			http.Error(w, fmt.Sprintf("%q not found", key), http.StatusNotFound)
			return
		}
		fmt.Fprint(w, value)
	}
}

// handleSet returns a handler that retrieves key-value string pairs from the URI and stores their values in
// sb. Will both store empty strings as values, and clobber values for existing keys. One request can set
// several keys; each query parameter's name is the stored key, and its value is stored as the value. If
// multiple values are specified on a single key, the first value wins.
//
//	/set?somekey=somevalue
//
// Sets the value of "somekey" to "somevalue", regardless of whether "somekey" currently exists in sb.
//
// Note that /set accepts GET here because the spec only specifies a URL. An HTTP GET is supposed to be safe
// in addition to being idempotent, but this is where we find ourselves.
func handleSet(sb *stringabase.Stringabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for key, values := range query {
			sb.Set(key, values[0])
		}
	}
}
