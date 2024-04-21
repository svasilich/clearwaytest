package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

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

	dataServer := dataserverapp.NewDataServerApp(repo, repo, auth.HasherMD5Hex)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth", dataServer.Auth)
	mux.HandleFunc("/api/upload-asset/", dataServer.Upload)
	mux.HandleFunc("/api/asset/", dataServer.Download)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), mux))
}
