// Stringabase is a barebones, concurrency-safe, string key-value store served over a simple HTTP API on TCP
// port 4000. It's unfussy about values, it clobbers, all the defaults are baked in, and it's strictly
// in-memory.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ianeff/stringabase/internal/server"
	"github.com/ianeff/stringabase/internal/stringabase"
)

func main() {
	addr := "localhost:4000"
	sb := stringabase.NewStringabase()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := server.ListenAndServe(ctx, addr, sb)
	if err != nil {
		log.Fatal(err)
	}
}
