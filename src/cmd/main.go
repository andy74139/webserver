package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andy74139/webserver/src/app"
)

// entrypoint
//
//go:generate swagger generate spec -o swagger.json
func main() {
	ctx := context.Background()
	a := app.New()
	go func() {
		if err := a.Start(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to start app: %s\n", err)
		}
	}()
	gracefulShutdown(ctx, a)
}

func gracefulShutdown(ctx context.Context, a app.App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.Stop(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed stop app: %v\n", err)
		}
	}()
	wg.Wait()
}
