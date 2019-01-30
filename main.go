// Package main is the core of the server - it handles the main routing, database initalization, etc.
package main

import (
	"database/sql"
	"goshare/config"
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
	config.Initialize()
	db, e := sql.Open("sqlite3", config.DBPath)
	if e != nil {
		log.Panicln(e)
	}
	defer db.Close()
	storage.Initialize(db)
	r.PathPrefix("/").HandlerFunc(handlers.Default).Methods("GET")
	r.HandleFunc("/upload", handlers.Upload).Methods("POST")
	log.Println("Server is ready.")
	e = http.ListenAndServe("127.0.0.1:"+config.Port, r)
	if e != nil {
		log.Panicln(e)
	}
}
