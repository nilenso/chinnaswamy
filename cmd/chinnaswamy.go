package main

import (
	"context"
	"fmt"
	"nilenso.com/chinnaswamy"
	"nilenso.com/chinnaswamy/log"
	"os"
	"os/signal"
)

func main() {
	cfg := map[string]string{"listenAddress": ":8080"}
	ctx := context.Background()

	err := log.InitLogger()
	if err != nil {
		_ = fmt.Errorf("error: Could not start the logger")
		os.Exit(2)
	}

	defer func() {
		err := log.SyncLogs()
		if err != nil {
			_ = fmt.Errorf("error: Could not sync the log")
			os.Exit(3)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	serverDone := make(chan struct{}, 1)
	go chinnaswamy.Serve(ctx, cfg, serverDone)
	log.Infow("Server started", "listenAddress", cfg["listenAddress"])

	select {
	case <-serverDone:
		log.Infow("Server exited. Exiting program")
		os.Exit(0)
	case <-sigint:
		log.Infow("Interrupt received. Exiting")
		os.Exit(0)
	}
}
