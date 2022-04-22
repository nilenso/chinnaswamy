package main

import (
	"fmt"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/shorten"
	"os"
	"time"
)

func main() {
	if err := log.InitLogger(); err != nil {
		_ = fmt.Errorf("error: Could not start the logger: %s", err)
		os.Exit(2)
	}

	client := shorten.NewClient("localhost:8090")
	err := client.Connect()
	defer client.Close()
	if err != nil {
		log.Errorw("Could not connect to gRPC server",
			"errorMessage", err,
		)
		os.Exit(1)
	}
	urlMapping, err := client.ShortenUrl("https://nilenso.com", 5*time.Minute)
	if err != nil {
		log.Errorw("Could not shorten URL",
			"errorMessage", err,
		)
		os.Exit(2)
	}
	log.Infow("Created short URL", "urlMapping", urlMapping)
}
