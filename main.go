package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jtcressy/vmware-work-sample/pkg/synthetics"
)

var (
    options synthetics.Options
)

func main() {
    defaultFetchInterval, _ := time.ParseDuration("30s")
    flag.DurationVar(
        &options.FetchInterval,
        "fetch-interval",
        defaultFetchInterval,
        "How often to ping URL endpoints (default: 30s)",
    )
    flag.Var(
        &options.TestUrls,
        "test-url",
        "A URL to query for uptime statistics. Use multiple times to query multple URL's in parallel",
    )
    flag.StringVar(
        &options.BindAddress,
        "bind-addr",
        ":8080",
        "The address the webserver binds to.",
    )
    flag.Parse()

    pinger, err := synthetics.NewPinger(options)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    ctx, cancel := context.WithCancel(context.Background())
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cancel()
        <-c
        os.Exit(1)
    }()
    if err := pinger.Start(ctx); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	// http.Handle("/metrics", promhttp.Handler())
    // http.Handle("/status", http.DefaultServeMux)
	// http.ListenAndServe(":8080", nil)
}