package stringabase_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ianeff/stringabase/internal/stringabase"
)

func TestSet_SetsTheValueOfSomeKeyInTheMapAndGetGetsIt(t *testing.T) {
	t.Parallel()
	db := newTestStore(t)
	want := "1"

	db.Set("a", "1")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}
	if got != want {
		t.Errorf("wrong value after Set (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestSet_LastOneWinsAndClobbers(t *testing.T) {
	t.Parallel()
	db := newTestStore(t)
	want := "2"

	db.Set("a", "1")
	db.Set("a", "2")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}

	if got != want {
		t.Errorf("second Set didn't clobber the first (-want +got):\n%s", cmp.Diff(want, got))
	}
}

func TestGet_DistinguishesAnEmptyValueFromAMissingKey(t *testing.T) {
	t.Parallel()
	db := newTestStore(t)
	wantA := ""
	wantB := ""

	db.Set("a", "")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}
	if got != wantA {
		t.Errorf("wrong value for a key set to empty (-want +got):\n%s", cmp.Diff(wantA, got))
	}

	got, ok = db.Get("b")
	if ok {
		t.Fatal("missing key found")
	}
	if got != wantB {
		t.Errorf("missing key returned a value (-want +got):\n%s", cmp.Diff(wantB, got))
	}
}

// TestSet_IsRaceFree exercises the store with a concurrent reader and writer so the race detector can catch
// unsynchronized access. The channel is established so we wait on the writer to finish before exiting the
// test.
func TestSet_IsRaceFree(t *testing.T) {
	t.Parallel()
	db := newTestStore(t)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 100 {
			db.Set("a", "1")
		}
	}()
	for range 100 {
		db.Get("a")
	}
	// Let the writer finish before we return.
	<-done
}

func newTestStore(t *testing.T) *stringabase.Store {
	t.Helper()
	return stringabase.NewStore()
}
