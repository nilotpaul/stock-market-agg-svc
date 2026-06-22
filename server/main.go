package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var (
		dbHost     = os.Getenv("DB_HOST")
		dbKeyspace = os.Getenv("DB_KEYSPACE")
		listenAddr = os.Getenv("API_LISTEN_ADDR")
	)
	if len(dbHost) == 0 {
		log.Fatal("db host environment variable not provided")
	}
	if len(dbKeyspace) == 0 {
		log.Fatal("db keyspace environment variable not provided")
	}
	if len(listenAddr) == 0 {
		log.Fatal("api listen addr environment variable not provided")
	}

	sess, err := OpenDBSession(dbHost, dbKeyspace)
	if err != nil {
		log.Fatal("db error:", err)
	}
	defer sess.Close()

	s := NewAPIServer(listenAddr, sess)
	if err := s.Start(); err != nil {
		log.Println("api serve error:", err)
	}
}
