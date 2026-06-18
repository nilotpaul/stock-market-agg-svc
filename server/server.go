package main

import (
	"fmt"
	"net/http"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/nilotpaul/stock-market-agg-svc/server/repository"
	"github.com/nilotpaul/stock-market-agg-svc/server/service"
)

type APIServer struct {
	listenAddr string
	session    *gocql.Session
}

func NewAPIServer(listenAddr string, sess *gocql.Session) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		session:    sess,
	}
}

func (s *APIServer) Start() error {
	s.registerAPIRoutes()

	fmt.Printf("started at %s\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *APIServer) registerAPIRoutes() {
	dbRepo := repository.NewCassandra(s.session)
	svc := service.NewCandleService(dbRepo)

	ch := NewCandleHandler(svc)

	http.HandleFunc("GET /api/v1/health", handler(handleGetHealth))
	http.HandleFunc("GET /api/v1/candles", handler(ch.HandleGetCandles))
}
