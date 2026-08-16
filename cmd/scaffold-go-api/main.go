// Command scaffold-go-api is the composition root for the scaffold's example
// service: it wires the in-memory widget store to the handlers and starts
// the server. A project stamped from this template replaces the store with
// a real one and grows internal/ from there — see CLAUDE.md.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	scaffoldgoapi "github.com/alrayyes/scaffold-go-api"
)

// version is stamped in at build time by goreleaser, from the tag. "dev" is
// what a plain `go build` reports, which is the honest answer for a binary
// built off an unknown tree.
var version = "dev"

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	widgets := map[string]scaffoldgoapi.Widget{
		"hammer": {ID: "hammer", Name: "Claw hammer"},
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           scaffoldgoapi.NewMux(widgets),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown", "error", err)
		}
	}()

	slog.Info("starting", "version", version, "addr", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
