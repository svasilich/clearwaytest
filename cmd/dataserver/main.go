package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/svasilich/clearwaytest/internal/application/auth"
	"github.com/svasilich/clearwaytest/internal/application/dataserverapp"
	"github.com/svasilich/clearwaytest/internal/repository/cwrepo"
)

func main() {
	var port int
	var dbConnectionString string
	flag.IntVar(&port, "port", 8080, "port for server")
	flag.StringVar(&dbConnectionString, "dbcon", "postgres://postgres:postgres@localhost:5432/postgres?pool_max_conns=10", "db connection string")
	flag.Parse()

	repo := cwrepo.NewRepository(dbConnectionString)
	if err := repo.Connect(context.Background()); err != nil {
		log.Fatalf("can't connect to repository: %s", err.Error())
	}
	defer repo.Close()

	dataServer := dataserverapp.NewDataServerApp(repo, repo, auth.HasherMD5Hex, repo, repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth", loggerMiddleware(dataServer.Auth))
	mux.HandleFunc("/api/upload-asset/", loggerMiddleware(dataServer.Upload))
	mux.HandleFunc("/api/asset/", loggerMiddleware(dataServer.Download))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), mux))
}

func loggerMiddleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := uuid.New() // Всё же использую здесь внешний пакет, чтобы удобно было генерировать reqid'ы.
		start := time.Now()
		log.Printf("start handle %s with reqid %s", r.RequestURI, reqid.String())

		handler(w, r)

		log.Printf("request with reqid %s was handled. Time is %v", reqid.String(), time.Since(start))
	})
}
