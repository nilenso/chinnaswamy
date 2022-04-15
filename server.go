package chinnaswamy

import (
	"context"
	"encoding/json"
	"net/http"
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

func Serve(ctx context.Context, cfg map[string]string, done chan struct{}) {
	srv := &http.Server{
		Addr:         cfg["listenAddress"],
		WriteTimeout: 30 * time.Millisecond,
		ReadTimeout:  30 * time.Millisecond,
		IdleTimeout:  1 * time.Second,
		Handler:      defaultMux(),
	}
	log.Infow("Starting server",
		"listenAddress", cfg["listenAddress"],
	)
	err := srv.ListenAndServe()
	if err != nil {
		log.Errorw("Could not start server",
			"listenAddress", cfg["listenAddress"],
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
