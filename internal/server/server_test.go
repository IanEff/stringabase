package server_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ianeff/stringabase/internal/server"
	"github.com/ianeff/stringabase/internal/stringabase"
)

func TestGet_ReturnsAValueStoredByASet(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := "1"

	fetch(t, ts, "/set?a=1")
	status, got := fetch(t, ts, "/get?key=a")

	if status != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, status)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSet_ClobbersValueFromEarlierKey(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := "2"

	fetch(t, ts, "/set?a=1")
	fetch(t, ts, "/set?a=2")
	_, got := fetch(t, ts, "/get?key=a")

	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGet_RespondsWithNotFoundForBogusKey(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)

	status, _ := fetch(t, ts, "/get?key=a")

	if status != http.StatusNotFound {
		t.Errorf("want status %d, got %d", http.StatusNotFound, status)
	}
}

func TestGet_DistinguishesAnEmptyValueFromAMissingKey(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)

	fetch(t, ts, "/set?a=")
	status, got := fetch(t, ts, "/get?key=a")

	if status != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, status)
	}
	if got != "" {
		t.Error(cmp.Diff("", got))
	}
}

func TestListenAndServe_ReturnsWhenContextIsCanceled(t *testing.T) {
	t.Parallel()
	// Hand the HTTP server a canceled context before it starts listening so it returns immediately.
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	errc := make(chan error, 1)
	go func() {
		errc <- server.ListenAndServe(ctx, randomLocalAddr(t), stringabase.NewStringabase())
	}()
	select {
	case err := <-errc:
		if err != nil {
			t.Errorf("want nil error after cancel, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("ListenAndServe did not return after ctx was canceled")
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewTestServer(t, server.NewMuxer(stringabase.NewStringabase()))
}

func fetch(t *testing.T, ts *httptest.Server, path string) (int, string) {
	t.Helper()
	resp, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(body)
}

func randomLocalAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}
