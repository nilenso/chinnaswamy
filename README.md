# Chinnaswamy
The 'Scalable' URL shortener

## Running tests

To run end to end tests:
```
CHINNASWAMY_TEST_HOST="<host:port>" go test -count=1 ./integration_test
```

## Running locally

Use docker compose:
```
docker-compose -f docker-compose.dev.yaml up -d
```

NOTE: You must ensure that a keyspace named `chinnaswamy` exists in the Cassandra cluster.

Eg, for local setup:
```cassandraql
CREATE KEYSPACE chinnaswamy WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
};
```

## Running migrations

Use the `migrate` docker image, like so:

```
docker run -v <path/to/repo>/db/migrations:/migrations \
           --network <docker-network> \
           migrate/migrate -path=/migrations/ -database cassandra://<cassandra host>/chinnaswamy up
```

For the local development setup started with:

```
docker-compose -f docker-compose.dev.yaml up -d
```

This would look like:
```
docker run -v ~/chinnaswamy/db/migrations:/migrations \
            --network chinnaswamy_chinnaswamy-devnet \
            migrate/migrate -path=/migrations/ -database cassandra://cassandra:9042/system up
```