// Package main is the core of the server - it handles the main routing, database initalization, etc.
package main

import (
	"database/sql"
	"goshare/handlers"
	"goshare/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type test struct {
	key   string
	value string
}

func main() {
	var e error
	r := mux.NewRouter()
	db, e := sql.Open("sqlite3", "./test.db")
	if e != nil {
		log.Panicln(e)
	}
	storage.Initialize(db)
	r.PathPrefix("/").HandlerFunc(handlers.Default).Methods("GET")
	r.HandleFunc("/upload", handlers.Upload).Methods("POST")
	e = http.ListenAndServe("127.0.0.1:8080", r)
	if e != nil {
		log.Panicln(e)
	}
	defer db.Close()
}
