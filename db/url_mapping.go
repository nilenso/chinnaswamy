package db

import (
	"context"
	"errors"
	"math"
	"nilenso.com/chinnaswamy/log"
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

func (s *Session) QueryUrlMapping(ctx context.Context, shortUrl string) (*urlmapping.UrlMapping, error) {
	urlMapping := &urlmapping.UrlMapping{}
	var valid bool
	query := `SELECT short_url, long_url, valid, TTL(valid) from url_mappings
              WHERE short_url = ?`
	err := s.Query(query, shortUrl).WithContext(ctx).Scan(
		&urlMapping.ShortUrl, &urlMapping.LongUrl, &valid, &urlMapping.TTL,
	)
	if err != nil {
		return nil, err
	}
	if !valid {
		log.Infow("Client tried to access invalid link")
		return nil, errors.New("tried to access invalid link")
	}
	return urlMapping, nil
}
