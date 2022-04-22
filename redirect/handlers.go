package redirect

import (
	"context"
	"net/http"
	"nilenso.com/chinnaswamy/log"
	"time"
)

func (rs *RedirectionService) queryUrl(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	shortUrl := req.URL.Path[1:]
	mapping, err := rs.dbSession.QueryUrlMapping(ctx, shortUrl)
	if err != nil {
		log.Errorw("Could not resolve URL mapping",
			"errorMessage", err,
		)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", mapping.LongUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
