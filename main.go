package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rasyncan/cmd"
	"syscall"
	"time"
)

func main() {
	fmt.Println("raSyncan")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	config := cmd.SyncConfig{}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go signalHandler(sigChan, config, ctx, cancel)

	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	start := time.Now()
	fmt.Println("start Time: ", start)

	if err := cmd.Execute(ctx, config, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Println("TTR: ", time.Now().Sub(start).Milliseconds(), "ms")
}

func signalHandler(sigChan chan os.Signal, config cmd.SyncConfig, ctx context.Context, cancel context.CancelFunc) {
	ctx = context.WithoutCancel(ctx)
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Println("We're done here received: ", sig)
				cancel()
				os.Exit(1)
			case syscall.SIGHUP:
				config.Initialize()
				log.Println("Reloading config, received: ", sig)
			}
		case <-ctx.Done():
			os.Exit(1)
		}
	}
}

//core functionality.
//File Synchronization: Implement the core functionality to compare files and directories
//between source and destination locations and synchronize them accordingly.

//compare files...
//same name
//same size,
//ue rsync algo to check if the file content is inherently the same.

//  rsync algo
