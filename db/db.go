package db

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/gocql/gocql"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/log"
	"time"
)

type Session struct {
	*gocql.Session
}

const keyspace = "chinnaswamy"

func initOnce() (*Session, error) {
	var session Session
	var err error
	log.Infow("Starting database on the given addresses",
		"addresses", config.DatabaseAddresses(),
	)
	cluster := gocql.NewCluster(config.DatabaseAddresses()...)
	cluster.Keyspace = keyspace
	session.Session, err = cluster.CreateSession()
	return &session, err
}

func retry[V comparable](fn func() (V, error), b backoff.BackOff) (V, error) {
	var val *V
	var zeroVal V
	err := backoff.Retry(func() error {
		ret, err := fn()
		if err != nil {
			return err
		}

		if ret != zeroVal {
			val = &ret
		}
		return nil
	}, b)

	if err != nil {
		return zeroVal, err
	}

	if val != nil {
		return *val, nil
	}

	return zeroVal, nil
}

func Init() (*Session, error) {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 3 * time.Minute
	expBackoff.Reset()
	return retry(initOnce, expBackoff)
}

func (s *Session) IsAvailable() bool {
	return s.Query("SELECT uuid() FROM system.local").Exec() == nil
}
