// Package server exposes [stringabase.Stringabase] as an HTTP API. The endpoints are documented in [ListenAndServe].
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ianeff/stringabase/internal/stringabase"
)

// ListenAndServe serves sb out over HTTP at addr on endpoints /get and /set, described in [NewMuxer], until ctx is canceled. sb is an intialized [stringabase.Stringabase]. addr takes any form net/http accepts ("localhost:3000").
// When ctx is canceled, ListenAndServe drops in-flight requests and returns.  If it cannot serve, it immediately returns a non-nil error.
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
// It collects and exposes the following endpoints:
//
// /get retrieves a value from the Stringabase by key.  Returns an HTTP StatusOK and the value if the key is present, even if that value is an empty string.  Returns a 404 if the key is not found.
//
//	/set sets a key-value string pair in the Stringabase.  Will both accept empty strings for values, and clobber the existing value of any existing key.  One request can set several keys.  First value wins.
func NewMuxer(sb *stringabase.Stringabase) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", handleGet(sb))
	mux.HandleFunc("/set", handleSet(sb))
	return mux
}

// handleGet retrieves a key from the URI and attempts to retrieve its value from sb.
func handleGet(sb *stringabase.Stringabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value, ok := sb.Get(key)
		if !ok {
			http.Error(w, fmt.Sprintf("%q not found", key), http.StatusNotFound)
			return
		}
		fmt.Fprint(w, value)
	}
}

// handleSet retrieves key-value string pairs from the URI and stores their values in sb. Will both store empty strings as values, and clobber values for existing keys.  One request can set several keys.  If multiple values are specified on a single key, the first value wins.
func handleSet(sb *stringabase.Stringabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.URL.Query() {
			sb.Set(key, values[0])
		}
	}
}
