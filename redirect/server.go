package redirect

import (
	"context"
	"encoding/json"
	"net/http"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/db"
	"nilenso.com/chinnaswamy/log"
	"nilenso.com/chinnaswamy/urlmapping"
	"time"
)

const (
	Available   = "AVAILABLE"
	Unavailable = "UNAVAILABLE"
)

type UrlMappingStore interface {
	QueryUrlMapping(ctx context.Context, shortUrl string) (urlmapping.UrlMapping, error)
}

type RedirectionService struct {
	dbSession *db.Session
}

func NewRedirectionService(dbSession *db.Session) *RedirectionService {
	return &RedirectionService{dbSession: dbSession}
}

func (rs *RedirectionService) isDatabaseAvailable() string {
	if rs.dbSession.IsAvailable() {
		return Available
	} else {
		return Unavailable
	}
}

func (rs *RedirectionService) defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(writer http.ResponseWriter, req *http.Request) {
		responseMap := map[string]string{"redirect": Available}
		responseMap["database"] = rs.isDatabaseAvailable()

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(responseMap)
		if err != nil {
			http.Error(writer, "{\"redirect\": \"UNAVAILABLE\"}", http.StatusInternalServerError)
			return
		}
	})
	return mux
}

func (rs *RedirectionService) Serve(ctx context.Context, done chan struct{}) {
	srv := &http.Server{
		Addr:         config.RedirectListenAddress(),
		WriteTimeout: config.WriteTimeout(),
		ReadTimeout:  config.ReadTimeout(),
		IdleTimeout:  config.IdleTimeout(),
		Handler:      rs.defaultMux(),
	}
	log.Infow("Starting redirection server",
		"listenAddress", config.RedirectListenAddress(),
	)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Errorw("Could not shut down redirection server gracefully")
			return
		}
		log.Infow("Redirection server has been shut down gracefully")
	}()

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorw("Could not start redirection server",
			"listenAddress", config.RedirectListenAddress(),
			"errorMessage", err,
		)
		done <- struct{}{}
		return
	}

	done <- struct{}{}
	return
}
