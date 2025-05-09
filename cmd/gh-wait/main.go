package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nobe4/gh-wait/internal/checker"
	"github.com/nobe4/gh-wait/internal/github"
	"github.com/nobe4/gh-wait/internal/looper"
	"github.com/nobe4/gh-wait/internal/looper/callback"
	"github.com/nobe4/gh-wait/internal/version"
)

const (
	minDelay = 10
)

func main() {
	clearScreen := flag.Bool("clear", false, "clear the screen")
	condition := flag.String("condition", "closed", "condition to wait for")
	delay := flag.Int("delay", minDelay, "delay in seconds between each checks")
	versionFlag := flag.Bool("version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [FLAGS]\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Conditions:\n%s", "todo")
	}
	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "version: %s\n", version.String())
		os.Exit(0)
	}

	args := flag.Args()

	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	if err := run(*delay, *clearScreen, *condition, url); err != nil {
		panic(err)
	}
}

//revive:disable:cognitive-complexity // TODO: refactor.
func run(delay int, clr bool, condition, url string) error {
	p, err := github.ParsePull(url)
	if err != nil {
		panic(err)
	}

	c := checker.Get(condition)
	if c == nil {
		return fmt.Errorf("unknown condition: %s", condition)
	}

	var loop looper.Looper

	loop = callback.New(func() {
		if clr {
			clearScreen()
		}

		checked, msg, err := c.Check(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error checking: %v", err)

			return
		}

		fmt.Fprint(os.Stdout, msg+"\n")

		if checked {
			loop.Stop()
		}
	}, time.Second*time.Duration(delay))

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := loop.Loop(ctx); err != nil {
		return fmt.Errorf("error in loop: %w", err)
	}

	return nil
}

func clearScreen() {
	// TODO: clear screen
}
