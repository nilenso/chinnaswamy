#!/bin/sh

CASSANDRA_KEYSPACE="chinnaswamy"

echo "Creating keyspace..."
if test -n "$CASSANDRA_KEYSPACE"; then
  echo "Attempting to create keyspace: $CASSANDRA_KEYSPACE"
  CQL="CREATE KEYSPACE $CASSANDRA_KEYSPACE WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"
  until echo "$CQL" | cqlsh; do
    echo "cqlsh: Cassandra is unavailable - retry later"
    sleep 2
  done &
fi

echo "Going to the docker entrypoint..."
exec /usr/local/bin/docker-entrypoint.sh "$@"