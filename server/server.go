package server

import (
	"context"
	"encoding/json"
	"net/http"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/db"
	"nilenso.com/chinnaswamy/log"
	"time"
)

const (
	Available   = "AVAILABLE"
	Unavailable = "UNAVAILABLE"
)

func isDatabaseAvailable() string {
	if db.IsAvailable() {
		return Available
	} else {
		return Unavailable
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(writer http.ResponseWriter, req *http.Request) {
		responseMap := map[string]string{"server": Available}
		responseMap["database"] = isDatabaseAvailable()

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(responseMap)
		if err != nil {
			http.Error(writer, "{\"server\": \"UNAVAILABLE\"}", http.StatusInternalServerError)
			return
		}
	})
	return mux
}

func Serve(ctx context.Context, done chan struct{}) {
	srv := &http.Server{
		Addr:         config.ListenAddress(),
		WriteTimeout: config.WriteTimeout(),
		ReadTimeout:  config.ReadTimeout(),
		IdleTimeout:  config.IdleTimeout(),
		Handler:      defaultMux(),
	}
	log.Infow("Starting server",
		"listenAddress", config.ListenAddress(),
	)
	err := srv.ListenAndServe()
	if err != nil {
		log.Errorw("Could not start server",
			"listenAddress", config.ListenAddress(),
		)
		done <- struct{}{}
		return
	}

	go func(srv *http.Server) {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Errorw("Could not shut down server gracefully")
			return
		}
		log.Infow("Server has been shut down gracefully")
	}(srv)

	done <- struct{}{}
	return
}
