package main

import (
	"avito_task/httpapi"
	"avito_task/service/postgersIMPL"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	user     = "root"
	password = "root"
	host     = "postgres"
	port     = "5432"
	dbname   = "taskdb"
)

func dsn() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, dbname)
}

func main() {
	srv := httpapi.NewServer(postgersIMPL.NewManager(dsn()))
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
