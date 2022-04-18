package main

import (
	"context"
	"fmt"
	"nilenso.com/chinnaswamy"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/log"
	"os"
	"os/signal"
)

func main() {
	logErr := log.InitLogger()
	if logErr != nil {
		_ = fmt.Errorf("error: Could not start the logger")
		os.Exit(2)
	}

	cfgErr := config.Init()
	if cfgErr != nil {
		log.Errorf("Error loading config: %s. Exiting", cfgErr)
		os.Exit(1)
	}
	ctx := context.Background()

	defer func() {
		syncErr := log.SyncLogs()
		if syncErr != nil {
			_ = fmt.Errorf("error: Could not sync the log")
			os.Exit(3)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	serverDone := make(chan struct{}, 1)
	go chinnaswamy.Serve(ctx, serverDone)
	log.Infow("Server started", "listenAddress", config.ListenAddress())

	select {
	case <-serverDone:
		log.Infow("Server exited. Exiting program")
		os.Exit(0)
	case <-sigint:
		log.Infow("Interrupt received. Exiting")
		os.Exit(0)
	}
}
