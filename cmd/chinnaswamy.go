package main

import (
	"context"
	"fmt"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/db"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/redirect"
	"nilenso.com/chinnaswamy/shorten"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := log.InitLogger(); err != nil {
		_ = fmt.Errorf("error: Could not start the logger: %s", err)
		os.Exit(2)
	}

	if err := config.Init(); err != nil {
		log.Errorw("Error loading config",
			"error", err,
		)
		os.Exit(1)
	}

	dbSession, dbErr := db.Init()
	if dbErr != nil {
		log.Errorw("Error connecting to the db",
			"error", dbErr,
		)
		os.Exit(4)
	}

	redirectionService := redirect.NewRedirectionService(dbSession)
	shorteningService := shorten.NewShorteningService(dbSession)

	defer func() {
		_ = log.SyncLogs()
		// When logging to terminal, sleeping for a second gives enough
		// time to flush
		time.Sleep(1 * time.Second)
		/* Commented out because: https://github.com/uber-go/zap/issues/880
		syncErr := log.SyncLogs()
		if syncErr != nil {
			_ = fmt.Errorf("error: Could not sync the log: %s", syncErr)
			os.Exit(3)
		}
		*/
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	redirectCtx, cancelRedirectionSvc := context.WithCancel(context.Background())
	defer cancelRedirectionSvc()
	shortenServerDone := make(chan struct{}, 1)

	shortenCtx, cancelShortenSvc := context.WithCancel(context.Background())
	defer cancelShortenSvc()
	redirectionServerDone := make(chan struct{}, 1)

	go redirectionService.Serve(redirectCtx, redirectionServerDone)
	log.Infow("Redirection service started", "listenAddress", config.RedirectListenAddress())

	go shorteningService.Serve(shortenCtx, shortenServerDone)
	log.Infow("Shorten service started", "listenAddress", config.ShortenListenAddress())

	select {
	case <-shortenServerDone:
		log.Infow("Shorten server exited")
		break
	case <-redirectionServerDone:
		log.Infow("Redirection server exited")
		break
	case <-sigint:
		log.Infow("Interrupt received")
		break
	}
}
