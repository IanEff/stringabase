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
		t.Errorf("wrong value from /get after /set (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestGet_ReturnsStatusOKForAnEmptyValue(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := ""

	fetch(t, ts, "/set?a=")
	status, got := fetch(t, ts, "/get?key=a")

	if status != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, status)
	}
	if want != got {
		t.Errorf("wrong value for a key set to empty (-want +got):\n%s", cmp.Diff(want, got))
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

func TestGet_RespondsWithBadRequestWhenKeyParameterIsMissing(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)

	status, _ := fetch(t, ts, "/get")

	if status != http.StatusBadRequest {
		t.Errorf("want status %d, got %d", http.StatusBadRequest, status)
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
		t.Errorf("second /set didn't clobber the first (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestSet_StoresEveryKeyInAMultiKeyRequest(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	wantA := "1"
	wantB := "2"

	fetch(t, ts, "/set?a=1&b=2")
	_, gotA := fetch(t, ts, "/get?key=a")
	_, gotB := fetch(t, ts, "/get?key=b")

	if wantA != gotA {
		t.Errorf("wrong value for key a (-want +got):\n%s", cmp.Diff(wantA, gotA))
	}
	if wantB != gotB {
		t.Errorf("wrong value for key b (-want +got):\n%s", cmp.Diff(wantB, gotB))
	}
}

func TestSet_KeepsFirstValueWithRepeatedKey(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := "1"

	fetch(t, ts, "/set?a=1&a=2")
	_, got := fetch(t, ts, "/get?key=a")

	if want != got {
		t.Errorf("repeated key didn't keep its first value (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestSet_ReturnsAStatusOKWithAnEmptyBody(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := ""

	status, got := fetch(t, ts, "/set?a=1")

	if status != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, status)
	}
	if want != got {
		t.Errorf("/set returned a body (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestGet_ReturnsAValueStoredWithURLEncodedKeyAndValue(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	want := "a+b&c"

	fetch(t, ts, "/set?encoded%20mess=a%2Bb%26c")
	status, got := fetch(t, ts, "/get?key=encoded+mess")

	if status != http.StatusOK {
		t.Fatalf("want status %d, got %d", http.StatusOK, status)
	}
	if want != got {
		t.Errorf("wrong value for URL-encoded key and value (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestSet_RejectsAndDoesNotStoreOnBadQuery(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)

	status, _ := fetch(t, ts, "/set?a=1&b=100%")
	gotStatus, _ := fetch(t, ts, "/get?key=a")

	if status != http.StatusBadRequest {
		t.Errorf("want status %d, got %d", http.StatusBadRequest, status)
	}
	if gotStatus != http.StatusNotFound {
		t.Errorf("want status %d, got %d", http.StatusNotFound, gotStatus)
	}
}

func TestGet_RespondsWithBadRequestIfItCantDecodeAQuery(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)

	status, _ := fetch(t, ts, "/get?key=100%")

	if status != http.StatusBadRequest {
		t.Errorf("want status %d, got %d", http.StatusBadRequest, status)
	}
}

func TestListenAndServe_ReturnsNilWhenCanceledWhileServing(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	addr := randomLocalAddr(t)
	errc := make(chan error, 1)
	go func() {
		errc <- server.ListenAndServe(ctx, addr, stringabase.NewStore())
	}()
	waitUntilServing(t, addr)

	cancel()

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
	return httptest.NewTestServer(t, server.NewMuxer(stringabase.NewStore()))
}

func fetch(t *testing.T, ts *httptest.Server, path string) (int, string) {
	t.Helper()
	client := ts.Client()
	resp, err := client.Get(ts.URL + path)
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

func waitUntilServing(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		resp, err := http.Get("http://" + addr + "/get?key=a")
		if err == nil {
			resp.Body.Close()
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("server at %s never started: %v", addr, err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
