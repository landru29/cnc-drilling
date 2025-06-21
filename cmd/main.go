package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh

		cancel()
	}()

	mainCmd, err := mainCommand()
	if err != nil {
		panic(err)
	}

	if err := mainCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}
