package stringabase_test

import (
	"testing"

	"github.com/ianeff/stringabase/internal/stringabase"
)

func TestSet_SetsTheValueOfSomeKeyInTheMapAndGetGetsIt(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStringabase()
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
	db := stringabase.NewStringabase()
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

func TestSet_IsRaceFree(t *testing.T) {
	t.Parallel()
	db := stringabase.NewStringabase()
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
	<-done
}
