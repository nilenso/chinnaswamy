package db

import (
	"context"
	"math"
	"nilenso.com/chinnaswamy/urlmapping"
)

func (s *Session) StoreUrlMapping(ctx context.Context, urlMapping urlmapping.UrlMapping) error {
	insertQuery := `INSERT INTO url_mappings (short_url, long_url, created_at)
                    VALUES (?, ?, toTimestamp(now()))`
	insertErr := s.Query(insertQuery,
		urlMapping.ShortUrl, urlMapping.LongUrl,
	).WithContext(ctx).Exec()
	if insertErr != nil {
		return insertErr
	}
	updateQuery := `UPDATE url_mappings USING TTL ?
                    SET valid = ? WHERE short_url = ?`
	updateErr := s.Query(updateQuery,
		int(math.Round(urlMapping.TTL.Seconds())),
		urlMapping.TTL > 0, urlMapping.ShortUrl,
	).WithContext(ctx).Exec()
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (s *Session) QueryUrlMapping(ctx context.Context, shortUrl string) (urlmapping.UrlMapping, error) {
	urlMapping := urlmapping.UrlMapping{}
	err := s.Query("SELECT now() from system.local").WithContext(ctx).Scan(&urlMapping)
	return urlMapping, err
}
