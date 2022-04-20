package db

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/gocql/gocql"
	"nilenso.com/chinnaswamy/config"
	"nilenso.com/chinnaswamy/log"
	"time"
)

const keyspace = "chinnaswamy"

var cluster *gocql.ClusterConfig
var session *gocql.Session

func initOnce() error {
	var err error
	log.Infow("Starting database on the given addresses",
		"addresses", config.DatabaseAddresses(),
	)
	cluster = gocql.NewCluster(config.DatabaseAddresses()...)
	cluster.Keyspace = keyspace
	session, err = cluster.CreateSession()
	return err
}

func Init() error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 3 * time.Minute
	expBackoff.Reset()
	return backoff.Retry(initOnce, expBackoff)
}

func IsAvailable() bool {
	return session.Query("SELECT uuid() FROM system.local").Exec() == nil
}
