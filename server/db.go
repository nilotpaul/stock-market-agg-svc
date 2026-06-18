package main

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func OpenDBSession(host, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace

	return cluster.CreateSession()
}
