package main

import (
	"context"
	"fmt"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/db"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/server"
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
		log.Errorw("Error loading config",
			"error", cfgErr,
		)
		os.Exit(1)
	}
	ctx := context.Background()

	dbErr := db.Init()
	if dbErr != nil {
		log.Errorw("Error connecting to the db",
			"error", dbErr,
		)
		os.Exit(4)
	}

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
	go server.Serve(ctx, serverDone)
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
