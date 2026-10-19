package stringabase_test

import (
	"testing"

	"github.com/ianeff/stringabase/internal/stringabase"
)

func TestSet_SetsTheValueOfSomeKeyInTheMapAndGetGetsIt(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStore()
	db.Set("a", "1")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}
	if got != "1" {
		t.Fatalf("wanted %q, got %q", "1", got)
	}
}

func TestSet_LastOneWinsAndClobbers(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStore()
	db.Set("a", "1")
	db.Set("a", "2")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}
	if got != "2" {
		t.Fatalf("wanted %q, got %q", "2", got)
	}
}

func TestGet_DistinguishesAnEmptyValueFromAMissingKey(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStore()
	db.Set("a", "")
	got, ok := db.Get("a")
	if !ok {
		t.Fatal("key not found")
	}
	if got != "" {
		t.Fatalf("wanted %q, got %q", "", got)
	}
	got, ok = db.Get("b")
	if ok {
		t.Fatal("missing key found")
	}
	if got != "" {
		t.Fatalf("wanted %q, got %q", "", got)
	}
}

// TestSet_IsRaceFree exercises the store with a concurrent reader and writer so the race detector can catch
// unsynchronized access. The channel is established so we wait on the writer to finish before exiting the
// test.
func TestSet_IsRaceFree(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStore()
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
