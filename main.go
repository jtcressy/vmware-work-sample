package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jtcressy/vmware-work-sample/pkg/synthetics"
)

var (
	version string = "dev"
	options synthetics.Options
)

func main() {
	defaultFetchInterval := 5 * time.Second
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

	log.Println(os.Args[0], version, runtime.GOOS, runtime.GOARCH)

	pinger, err := synthetics.NewPinger(&options)
	if err != nil {
		log.Fatal(err)
	}
	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(pinger.GetContext().Context)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()
	if err := pinger.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
