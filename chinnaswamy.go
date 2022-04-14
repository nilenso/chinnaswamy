package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(writer http.ResponseWriter, req *http.Request) {
		responseMap := map[string]string{"status": "OK"}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(responseMap)
		if err != nil {
			http.Error(writer, "{\"status\": \"Error\"}", http.StatusInternalServerError)
			return
		}
	})
	return mux
}

func initLogger() *zap.SugaredLogger {
	logger, logInitErr := zap.NewDevelopment()
	if logInitErr != nil {
		_ = fmt.Errorf("error: Could not start the logger")
		os.Exit(2)
	}
	return logger.Sugar()
}

func main() {
	listenAddress := ":8080"

	log := initLogger()

	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			_ = fmt.Errorf("error: Could not sync the log")
			os.Exit(3)
		}
	}(log)

	err := http.ListenAndServe(listenAddress, defaultMux())
	if err != nil {
		log.Errorw("Server could not be started",
			"listenAddress", listenAddress,
		)
		return
	}

	log.Infow("Server started",
		"listenAddress", listenAddress,
	)
}
