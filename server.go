package chinnaswamy

import (
	"context"
	"encoding/json"
	"net/http"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/log"
	"time"
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

func Serve(ctx context.Context, cfg config.Config, done chan struct{}) {
	srv := &http.Server{
		Addr:         cfg.ListenAddress(),
		WriteTimeout: cfg.WriteTimeout(),
		ReadTimeout:  cfg.ReadTimeout(),
		IdleTimeout:  cfg.IdleTimeout(),
		Handler:      defaultMux(),
	}
	log.Infow("Starting server",
		"listenAddress", cfg.ListenAddress(),
	)
	err := srv.ListenAndServe()
	if err != nil {
		log.Errorw("Could not start server",
			"listenAddress", cfg.ListenAddress(),
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
