// Package server exposes [stringabase.Stringabase] as an HTTP API. The endpoints are documented in [ListenAndServe].
package server

// ListenAndServe serves sb out over HTTP at addr until ctx is canceled. sb is an intialized [stringabase.Stringabase]. addr takes any form net/http accepts ("localhost:3000").
// It services the following endpoints:
// 	/get retrieves a value from the Stringabase by key.  Returns an HTTP StatusOK and the value if the key is present.  Returns a 404 if the key is not found.
// 	/set sets a key-value string pair in the Stringabase.  Will both initialize empty values, and clobber existing values.
// When ctx is canceled, ListenAndServe drops in-flight requests and returns.  If it cannot serve, it immediately returns a non-nil error.
